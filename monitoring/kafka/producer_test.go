package kafka

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestNewProducer(t *testing.T) {
	fmt.Println(NewProducer(context.Background(), []string{"192.168.204.134:9092"}, 5))
}

func TestProducerMessage(t *testing.T) {
	var (
		producer *Producer
		err      error
		message  = make(chan []byte)
		timer    *time.Ticker
	)

	if producer, err = NewProducer(context.Background(), []string{"192.168.204.134:9092"}, 5); err != nil {
		t.Fatal(err)
	}

	go func() {
		if err = producer.SyncProducer("lucas", "luquitas", message); err != nil {
			t.Fatal(err)
		}
	}()

	timer = time.NewTicker(1 * time.Second)

	for {
		select {
		case _ = <-timer.C:
			fmt.Println("PRODUCER")
			message <- []byte("TESTE")
		}
	}
}
