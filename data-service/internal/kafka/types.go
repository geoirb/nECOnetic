package kafka

import (
	"context"
)

// Handler message from mq.
type Handler func(ctx context.Context, message []byte)

// Publish message to mq.
type Publish func(message []byte) error
