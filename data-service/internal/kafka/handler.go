package kafka

import (
	"context"

	"github.com/Shopify/sarama"
)

type consumerGroupHandler struct {
	ctx     context.Context
	handler Handler
}

func newConsumerGroupHandler(
	ctx context.Context,
	handler Handler,
) *consumerGroupHandler {
	return &consumerGroupHandler{
		ctx:     ctx,
		handler: handler,
	}
}

// ConsumeClaim processing data from kafka.
func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		h.handler(h.ctx, msg.Value)
	}
	return nil
}

// Setup ...
func (h *consumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup ...
func (h *consumerGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}
