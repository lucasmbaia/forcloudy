package main

import (
	"context"
	"fmt"
	"forcloudy/monitoring/monit"
)

func main() {
	fmt.Println(monit.Run(context.Background(), monit.Config{
		Running: 5,
		Topic:   "monitoring",
		Key:     "containers",
	}))
}
