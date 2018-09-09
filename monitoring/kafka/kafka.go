package kafka

import (
	"context"
	"github.com/Shopify/sarama"
	"log"
)

type Client struct {
	coon   sarama.SyncProducer
	config *sarama.Config
	ctx    context.Context
}

func NewProducer(ctx context.Context, hosts []string, retry int) (*Client, error) {
	var (
		client = &Client{}
		err    error
	)

	client.config = sarama.NewConfig()
	client.config.Producer.Retry.Max = retry
	client.ctx = ctx

	if client.coon, err = sarama.NewSyncProducer(hosts, nil); err != nil {
		return client, err
	}

	go func() {
		select {
		case _ = <-ctx.Done():
			if err = client.coon.Close(); err != nil {
				log.Println(err)
			}
		}
	}()

	return client, nil
}

func (c *Client) Producer(topic, key string, message <-chan []byte) error {
	var (
		err error
		pm  *sarama.ProducerMessage
	)

	for {
		select {
		case msg := <-message:
			pm = &sarama.ProducerMessage{
				Topic: topic,
				Key:   sarama.ByteEncoder([]byte(key)),
				Value: sarama.ByteEncoder(msg),
			}

			if _, _, err = c.coon.SendMessage(pm); err != nil {
				return err
			}
		case _ = <-c.ctx.Done():
			break
		}
	}

	return nil
}
