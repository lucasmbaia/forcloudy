package kafka

import (
	"context"
	"github.com/Shopify/sarama"
	"log"
)

type Consumer struct {
	coon sarama.Consumer
	ctx  context.Context
}

func NewConsumer(ctx context.Context, hosts []string) (*Consumer, error) {
	var (
		consumer = &Consumer{}
		err      error
	)

	consumer.ctx = ctx
	if consumer.coon, err = sarama.NewConsumer(hosts, nil); err != nil {
		return consumer, err
	}

	go func() {
		select {
		case _ = <-ctx.Done():
			if err = consumer.coon.Close(); err != nil {
				log.Println(err)
			}
		}
	}()

	return consumer, nil
}

func (c *Consumer) Consume(topic string, message chan<- []byte) error {
	var (
		consumer sarama.PartitionConsumer
		err      error
	)

	if consumer, err = c.coon.ConsumePartition(topic, 0, sarama.OffsetNewest); err != nil {
		return err
	}

	defer func() {
		if err = consumer.Close(); err != nil {
			log.Println(err)
		}
	}()

ConsumerLoop:
	for {
		select {
		case msg := <-consumer.Messages():
			message <- msg.Value
		case _ = <-c.ctx.Done():
			break ConsumerLoop
		}
	}

	return nil
}
