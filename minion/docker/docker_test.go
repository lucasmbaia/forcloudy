package docker

import (
	"context"
	"encoding/json"
	"log"
	"testing"
)

func TestDockerEvents(t *testing.T) {
	var (
		err   error
		errc  = make(chan error, 1)
		event = make(chan []byte)
	)

	go func() {
		if err = DockerEvents(context.Background(), event); err != nil {
			log.Panic(err)
		}
	}()

	for {
		select {
		case msg := <-event:
			var ev Events

			if err = json.Unmarshal(msg, &ev); err != nil {
				log.Panic(err)
			}

			log.Println(ev)
		case e := <-errc:
			log.Panic(e)
		}
	}
}
