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
    cn		[]string
    as		[]string
    app		[]string
  )

  if hostname, err = os.Hostname(); err != nil {
    errc <- err
  }

  if containers, err = docker.ListAllContainers(NO_FILTER); err != nil {
    errc <-err
  }

  go updateContainers(ctx, &containers)
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

	for _, container := range containers {
	  cn = strings.Split(container.Name, "_app-")
	  as = strings.Split(cn[1], "-")
	  app = strings.Join(as[:len(as)-1], "-")

	  if networks, err = utils.NetworkUtilization(container.PID, 1); err != nil {
	    log.Println(err)
	    continue
	  }


	}
	log.Println(containers)
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

func updateContainers(ctx context.Context, containers *[]docker.Containers) {
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

	*containers = append(*containers, c...)
      case DIE:
	for index, container := range *containers {
	  if container.Name == ev.Actor.Attributes.Name {
	    *containers = append((*containers)[:index], (*containers)[index + 1:]...)
	    break
	  }
	}
      }
    case e := <-errc:
      log.Println(e)
    }
  }
}
