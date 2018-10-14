package main

import (
  "forcloudy/monitoring/kafka"
  "forcloudy/monitoring/elastic-search"
  _metrics "forcloudy/monitoring/metrics"
  "encoding/json"
  "context"
  "log"
)

func main() {
  var (
    consumer  *kafka.Consumer
    client    *els.Client
    err	      error
    message   = make(chan []byte)
    ctx	      = context.Background()
  )

  if consumer, err = kafka.NewConsumer(ctx, []string{"172.16.95.183:9092"}); err != nil {
    log.Fatal(err)
  }

  if client, err = els.NewClient(ctx, "http://172.16.95.185:9200"); err != nil {
    log.Fatal(err)
  }

  go func() {
    if err = consumer.Consume("fbcm", message); err != nil {
      log.Fatal(err)
    }
  }()

  for {
    select {
    case msg := <-message:
      log.Println("MESSAGE: ", string(msg))
      var metrics _metrics.Customers

      if err = json.Unmarshal(msg, &metrics); err != nil {
	log.Fatal(err)
      }

      if err = client.Post(metrics.Name, "fbcm", metrics); err != nil {
	log.Fatal(err)
      }
    }
  }
}
