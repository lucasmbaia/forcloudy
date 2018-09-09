package kafka

import (
	"context"
	"fmt"
	"testing"
)

func TestNewConsumer(t *testing.T) {
	fmt.Println(NewConsumer(context.Background(), []string{"192.168.204.134:9092"}))
}

func TestConsume(t *testing.T) {
	var (
		consumer *Consumer
		err      error
		message  = make(chan []byte)
	)

	if consumer, err = NewConsumer(context.Background(), []string{"192.168.204.134:9092"}); err != nil {
		t.Fatal(err)
	}

	go func() {
		if err = consumer.Consume("monitoring", message); err != nil {
			t.Fatal(err)
		}
	}()

	for {
		select {
		case msg := <-message:
			fmt.Println(string(msg))
		}
	}
}
