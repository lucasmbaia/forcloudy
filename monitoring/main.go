package main

import (
  "forcloudy/monitoring/monit"
  "context"
)

func main() {
  monit.Run(context.Background(), 5)
}
