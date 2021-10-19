package kafka

import (
	"context"

	"github.com/Shopify/sarama"
)

type ConsumerGroup struct {
	ctx      context.Context
	consumer sarama.ConsumerGroup
}

func (c *ConsumerGroup) Handler(topics []string, h Handler) {
	go func() {
		defer c.consumer.Close()
		for {
			select {
			case <-c.ctx.Done():
				return
			default:
				if err := c.consumer.Consume(c.ctx, topics, newConsumerGroupHandler(c.ctx, h)); err != nil {
					return
				}
			}
		}
	}()
}
