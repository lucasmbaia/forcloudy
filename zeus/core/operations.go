package core

import (
  "github.com/lucasmbaia/go-xmpp/docker"
  "forcloudy/zeus/config"
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

func DeployAppication(d Deploy, append bool) error {
  var (
    image	    string
    keyApplication  string
    err		    error
    exists	    bool
  )

  image = fmt.Sprintf("%s_app-%s/image:%s", d.Customer, d.ApplicationName, d.ImageVersion)
  keyApplication = fmt.Sprintf("/%s/%s", d.Customer, d.ApplicationName)
  d.Image = image

  if err = config.EnvSingleton.EtcdConnection.Set(keyApplication, d); err != nil {
    return err
  }

  if len(minions) == 1 {
    if !append{
      for minion, _ := range minions {
	if exists, err = existsImage(image, minion); err != nil {
	  return err
	}

	if !exists {
	  if _, err = deploy(d, minion, d.ApplicationName, IMAGE_DEFAULT, true); err != nil {
	    return err
	  }

	  if err = generateImage(d.ApplicationName, d.ImageVersion, d.Build, minion); err != nil {
	    return err
	  }
	}
      }
    }
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
    return r.Elements.ID, nil
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
