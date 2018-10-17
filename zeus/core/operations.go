package core

import (
  "github.com/lucasmbaia/go-xmpp/docker"
  "forcloudy/zeus/config"
  "sync"
  "fmt"
)

const (
  IMAGE_DEFAULT	= "alpine"
)

type Deploy struct {
  Customer	  string	    `json:",omitempty"`
  ApplicationName string	    `json:",omitempty"`
  ImageVersion	  string	    `json:",omitempty"`
  Ports		  []Ports	    `json:",omitempty"`
  Args		  map[string]string `json:",omitempty"`
  Cpus		  string	    `json:",omitempty"`
  Memory	  string	    `json:",omitempty"`
  TotalContainers int		    `json:",omitempty"`
  Dns		  string	    `json:",omitempty"`
  Image		  string	    `json:",omitempty"`
  Build		  string	    `json:",omitempty"`
}

type Ports struct {
  Port	    int	    `json:",omitempty"`
  Protocol  string  `json:",omitempty"`
}

type MinionsCount struct {
  Name	string
  TotalContainers int
}

func DeployAppication(d Deploy, iterator int, append bool) error {
  var (
    image	    string
    keyApplication  string
    err		    error
    exists	    bool
    applicationName string
    wg		    sync.WaitGroup
    mainMinion	    string
    errc	    = make(chan error, 1)
    minionsCount    map[string]int
  )

  image = fmt.Sprintf("%s_app-%s/image:%s", d.Customer, d.ApplicationName, d.ImageVersion)
  applicationName = fmt.Sprintf("%s_app-%s", d.Customer, d.ApplicationName)
  keyApplication = fmt.Sprintf("/%s/%s", d.Customer, d.ApplicationName)
  d.Image = image

  if err = config.EnvSingleton.EtcdConnection.Set(keyApplication, d); err != nil {
    return err
  }

  if !append{
    for minion, _ := range minions {
      mainMinion = minion
    }

    if exists, err = existsImage(image, mainMinion); err != nil {
      return err
    }

    if !exists {
      if _, err = deploy(d, mainMinion, applicationName, IMAGE_DEFAULT, true); err != nil {
	return err
      }

      if err = generateImage(applicationName, d.ImageVersion, d.Build, mainMinion); err != nil {
	return err
      }

      if err = removeContainer(applicationName, mainMinion); err != nil {
	return err
      }
    }
  }

  if len(minions) == 1 {
    wg.Add(d.TotalContainers)

    for i := iterator; i <= d.TotalContainers; i++ {
      var containerName = fmt.Sprintf("%s_app-%s-%d", d.Customer, d.ApplicationName, i)
      go func(containerName string) {
	fmt.Println("MANDOU: ", containerName)
	if _, err = deploy(d, mainMinion, containerName, image, false); err != nil {
	  errc <- err
	}

	wg.Done()
      }(containerName)
    }

    go func() {
      select {
      case e := <-errc:
	err = e
      }
    }()

    wg.Wait()
  } else {
    if minionsCount, err = containersPerMinion(d.TotalContainers); err != nil {
      return err
    }

    wg.Add(d.TotalContainers)

    for key, value := range minionsCount {
      for i := 0; i < value; i++ {
	var containerName = fmt.Sprintf("%s_app-%s-%d", d.Customer, d.ApplicationName, iterator)
	go func(containerName, minion string) {
	  if _, err = deploy(d, minion, containerName, image, false); err != nil {
	    errc <- err
	  }

	  wg.Done()
	}(containerName, key)

	iterator++
      }
    }

    go func() {
      select {
      case e := <-errc:
	err = e
      }
    }()

    wg.Wait()
  }

  return nil
}

func existsImage(image, to string) (bool, error) {
  var (
    iq	      docker.IQ
    err	      error
    response  = make(chan ResponseIQ, 1)
  )

  if iq, err = docker.ExistsImage(docker.Image{
    From: config.EnvSingleton.XmppConnection.Jid,
    To:	  to,
    Name: image,
  }); err != nil {
    return false, err
  }

  mutex.Lock()
  responseIQ[iq.ID] = response
  mutex.Unlock()

  if err = config.EnvSingleton.XmppConnection.Send(iq); err != nil {
    return false, err
  }

  select {
  case r := <-response:
    if r.Elements.Name != EMPTY_STR {
      return true, nil
    }
  }

  return false, nil
}

func deploy(d Deploy, to, nameContainer, image string, imageCreate bool) (string, error) {
  var (
    iq	      docker.IQ
    err	      error
    response  = make(chan ResponseIQ, 1)
    ports     []docker.Ports
    args      []docker.Args
  )

  if len(d.Ports) > 0 {
    for _, port := range d.Ports {
      ports = append(ports, docker.Ports{Port: port.Port, Protocol: port.Protocol})
    }
  }

  if len(d.Args) > 0 {
    for key, value := range d.Args {
      args = append(args, docker.Args{Name: key, Value: value})
    }
  }

  if iq, err = docker.MasterDeploy(docker.Deploy{
    From:	      config.EnvSingleton.XmppConnection.Jid,
    To:		      to,
    Customer:	      d.Customer,
    ApplicationName:  d.ApplicationName,
    Name:	      nameContainer,
    Cpus:	      d.Cpus,
    Memory:	      d.Memory,
    Ports:	      ports,
    Args:	      args,
    Image:	      image,
    CreateImage:      imageCreate,
  }); err != nil {
    return EMPTY_STR, err
  }

  mutex.Lock()
  responseIQ[iq.ID] = response
  mutex.Unlock()

  if err = config.EnvSingleton.XmppConnection.Send(iq); err != nil {
    return EMPTY_STR, err
  }

  select {
  case r := <-response:
    return r.Elements.ID, r.Error
  }
}

func generateImage(image, version, build, to string) error {
  var (
    iq	      docker.IQ
    err	      error
    response  = make(chan ResponseIQ, 1)
  )

  if iq, err = docker.GenerateImage(docker.Image{
    From:	config.EnvSingleton.XmppConnection.Jid,
    To:		to,
    Name:	image,
    BuildName:	build,
    Tag:	version,
  }); err != nil {
    return err
  }

  mutex.Lock()
  responseIQ[iq.ID] = response
  mutex.Unlock()

  if err = config.EnvSingleton.XmppConnection.Send(iq); err != nil {
    return err
  }

  select {
  case _ = <-response:
    return nil
  }
}

func removeContainer(name, to string) error {
  var (
    iq	      docker.IQ
    err	      error
    response  = make(chan ResponseIQ, 1)
  )

  if iq, err = docker.RemoveContainer(docker.Action{
    From:	config.EnvSingleton.XmppConnection.Jid,
    To:		to,
    Container:	name,
  }); err != nil {
    return err
  }

  mutex.Lock()
  responseIQ[iq.ID] = response
  mutex.Unlock()

  if err = config.EnvSingleton.XmppConnection.Send(iq); err != nil {
    return err
  }

  select {
  case r := <-response:
    return r.Error
  }
}

func totalContainersMinion(to string) (int, error) {
  var (
    iq	      docker.IQ
    err	      error
    response  = make(chan ResponseIQ, 1)
  )

  if iq, err = docker.TotalContainers(config.EnvSingleton.XmppConnection.Jid, to); err != nil {
    return 0, err
  }

  mutex.Lock()
  responseIQ[iq.ID] = response
  mutex.Unlock()

  if err = config.EnvSingleton.XmppConnection.Send(iq); err != nil {
    return 0, err
  }

  select {
  case r := <-response:
    return r.Elements.TotalContainers, r.Error
  }
}

func containersPerMinion(totalContainers int) (map[string]int, error) {
  var (
    minionsCount      = make(map[int]MinionsCount)
    minionsContainers = make(map[string]int)
    err		      error
    total	      int
    totalMinions      = len(minions)
    count	      = 0
    differente	      bool
  )

  for minion, _ := range minions {
    if total, err = totalContainersMinion(minion); err != nil {
      return minionsContainers, err
    }

    if len(minionsCount) == 0 {
      minionsCount[count] = MinionsCount{Name: minion, TotalContainers: total}
    } else {
      for key, value := range minionsCount {
	if total < value.TotalContainers {
	  var mc = MinionsCount{Name: value.Name, TotalContainers: value.TotalContainers}
	  minionsCount[key] = MinionsCount{Name: minion, TotalContainers: total}
	  minionsCount[count] = mc
	} else {
	  minionsCount[count] = MinionsCount{Name: minion, TotalContainers: total}
	}
      }
    }

    count++
  }

  for i := 0; i < totalMinions; i++ {
    for j := i + 1; j < totalMinions; j++ {
      if minionsCount[i].TotalContainers != minionsCount[j].TotalContainers {
	differente = true
	break
      }
    }
  }

  if !differente {
    for _, value := range minionsCount {
      minionsContainers[value.Name] = totalContainers / totalMinions
    }

    if totalContainers % totalMinions != 0 {
      minionsContainers[minionsCount[0].Name] += totalContainers % totalMinions
    }
  }

  return minionsContainers, nil
}
