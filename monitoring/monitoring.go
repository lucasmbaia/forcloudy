package main

import (
  "context"
  "fmt"
  "forcloudy/monitoring/monit"
  "flag"
)

var (
  running = flag.Int("running", 5, "Time of check")
  kafka	  = flag.String("kafka", "172.16.95.183:9092", "URL of kafka")
)

func main() {
  flag.Parse()

  fmt.Println(monit.Run(context.Background(), monit.Config{
    Running:  *running,
    Topic:    "fbcm",
    Key:      "containers",
    Kafka:    []string{*kafka},
  }))
}
