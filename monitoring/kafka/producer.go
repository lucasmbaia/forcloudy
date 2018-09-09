package kafka

import (
	"context"
	"github.com/Shopify/sarama"
	"log"
)

type Producer struct {
	coon   sarama.SyncProducer
	config *sarama.Config
	ctx    context.Context
}

func NewProducer(ctx context.Context, hosts []string, retry int) (*Producer, error) {
	var (
		producer = &Producer{}
		err      error
	)

	producer.config = sarama.NewConfig()
	producer.config.Producer.Retry.Max = retry
	producer.ctx = ctx

	if producer.coon, err = sarama.NewSyncProducer(hosts, nil); err != nil {
		return producer, err
	}

	go func() {
		select {
		case _ = <-ctx.Done():
			if err = producer.coon.Close(); err != nil {
				log.Println(err)
			}
		}
	}()

	return producer, nil
}

func (p *Producer) SyncProducer(topic, key string, message <-chan []byte) error {
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

			if _, _, err = p.coon.SendMessage(pm); err != nil {
				return err
			}
		case _ = <-p.ctx.Done():
			break
		}
	}

	return nil
}
