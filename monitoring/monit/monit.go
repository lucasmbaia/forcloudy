package monit

import (
  "forcloudy/monitoring/docker"
  "forcloudy/monitoring/utils"
  _metrics "forcloudy/monitoring/metrics"
  "encoding/json"
  "context"
  "log"
  "time"
  "os"
  "sync"
  "strings"
)

const (
  NO_FILTER = ""
  START	    = "start"
  DIE	    = "die"
)

func Run(ctx context.Context, running int) error {
  var (
    containers	[]docker.Containers
    err		error
    errc	= make(chan error, 1)
    t		*time.Ticker
    hostname	string
    metrics	_metrics.Metrics
    customers	[]_metrics.Customers
    networks	[]_metrics.Networks
    net		= utils.NewNetwork()
    wg		sync.WaitGroup
    body	[]byte
  )

  if hostname, err = os.Hostname(); err != nil {
    errc <- err
  }

  if containers, err = docker.ListAllContainers(NO_FILTER); err != nil {
    errc <-err
  }

  go updateContainers(ctx, &containers, net)
  for _, container := range containers {
    if err = net.SetPid(container.PID); err != nil {
      errc <-err
    }
  }

  t = time.NewTicker(time.Duration(running) * time.Second)
  metrics.Hostname = hostname

  /*********
  adicionar sync.WaitGroup
  adcionar funcao para pegar o primeiro trafego de rede de cada container
  *********/

  go func() {
    for {
      select {
      case _ = <-t.C:
	customers = []_metrics.Customers{}
	wg.Add(len(containers))

	for _, container := range containers {
	  go func(container docker.Containers) {
	    var (
	      cn	[]string
	      as	[]string
	      app	string
	      containsC	bool
	      containsA	bool
	    )

	    cn = strings.Split(container.Name, "_app-")
	    as = strings.Split(cn[1], "-")
	    app = strings.Join(as[:len(as)-1], "-")

	    if networks, err = net.NetworkUtilization(container.PID, running); err != nil {
	      log.Println(err)
	    }

	    for indexC, customer := range customers {
	      if customer.Name == cn[0] {
		for indexA, application := range customers[indexC].Applications {
		  if app == application.Name {
		    customers[indexC].Applications[indexA].Containers = append(customers[indexC].Applications[indexA].Containers, _metrics.Containers{
		      ID:	container.ID,
		      Name:	container.Name,
		      Networks:	networks,
		    })

		    containsA = true
		    break
		  }
		}

		if !containsA {
		  customers[indexC].Applications = append(customers[indexC].Applications, _metrics.Applications{
		    Name: app,
		    Containers: []_metrics.Containers{{
		      ID:	    container.ID,
		      Name:	    container.Name,
		      Networks: networks,
		    }},
		  })
		}

		containsC = true
		break
	      }
	    }

	    if !containsC {
	      customers = append(customers, _metrics.Customers{
		Name: cn[0],
		Applications: []_metrics.Applications{{
		  Name: app,
		  Containers: []_metrics.Containers{{
		    ID:	    container.ID,
		    Name:	    container.Name,
		    Networks: networks,
		  }},
		}},
	      })
	    }

	    wg.Done()
	  }(container)
	}

	wg.Wait()
	metrics.Customers = customers

	if body, err = json.Marshal(metrics); err != nil {
	  log.Println(err)
	  break
	}

	log.Println(string(body))
      case _ = <-ctx.Done():
	log.Println("DONE")
      }
    }
  }()

  select {
  case err = <-errc:
    return err
  case _ = <-ctx.Done():
    return nil
  }
}

func updateContainers(ctx context.Context, containers *[]docker.Containers, net *utils.Utilization) {
  var (
    err	  error
    event = make(chan []byte)
    errc  = make(chan error, 1)
    c	  []docker.Containers
  )

  go func() {
    errc <- docker.DockerEvents(ctx, event)
  }()

  for {
    select {
    case msg := <-event:
      var ev docker.Events

      if err = json.Unmarshal(msg, &ev); err != nil {
	log.Println(err)
	continue
      }

      switch ev.Status {
      case START:
	if c, err = docker.ListAllContainers(ev.Actor.Attributes.Name); err != nil {
	  log.Println(err)
	  break
	}

	if err = net.SetPid(c[0].PID); err != nil {
	  log.Println(err)
	  break
	}

	*containers = append(*containers, c...)
      case DIE:
	for index, container := range *containers {
	  if container.Name == ev.Actor.Attributes.Name {
	    *containers = append((*containers)[:index], (*containers)[index + 1:]...)
	    net.DelPid(container.PID)
	    break
	  }
	}
      }
    case e := <-errc:
      log.Println(e)
    }
  }
}
