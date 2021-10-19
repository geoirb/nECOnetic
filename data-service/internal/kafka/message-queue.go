package kafka

import (
	"context"

	"github.com/Shopify/sarama"
)

// MessageQueue of kafka.
type MessageQueue struct {
	ctx       context.Context
	cancel    context.CancelFunc
	client    sarama.Client
	producer  sarama.SyncProducer
	consumers []*ConsumerGroup
}

// NewMessageQueue ...
func NewMessageQueue(
	addrs []string,
) (mq *MessageQueue, err error) {
	mq = &MessageQueue{}

	cfg := sarama.NewConfig()
	cfg.Version = sarama.MaxVersion
	cfg.Producer.Return.Successes = true

	if mq.client, err = sarama.NewClient(addrs, cfg); err != nil {
		return
	}
	if mq.producer, err = sarama.NewSyncProducerFromClient(mq.client); err != nil {
		return
	}
	mq.ctx, mq.cancel = context.WithCancel(context.Background())
	return
}

// NewConsumerGroup returns consumer for group id.
func (mq *MessageQueue) NewConsumerGroup(groupID string) (c *ConsumerGroup, err error) {
	c = &ConsumerGroup{
		ctx: mq.ctx,
	}
	if c.consumer, err = sarama.NewConsumerGroupFromClient(groupID, mq.client); err != nil {
		return
	}
	mq.consumers = append(mq.consumers, c)
	return
}

// NewPublisher returns publish func.
func (mq *MessageQueue) NewPublisher(topic string) Publish {
	return func(message []byte) (err error) {
		msg := &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.ByteEncoder(message),
		}
		_, _, err = mq.producer.SendMessage(msg)
		return
	}
}

// Shutdown consumers message queue.
func (mq *MessageQueue) Shutdown() {
	mq.cancel()
	mq.producer.Close()
}
