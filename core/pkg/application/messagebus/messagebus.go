package messagebus

import (
	"context"

	"github.com/mateusmacedo/golang-aether-sdk/core/pkg/application/message"
)

type Middleware func(next HandlerFunc) HandlerFunc

type HandlerFunc func(ctx context.Context, msg message.Message) error

type MessageBus interface {
    Publish(ctx context.Context, msg message.Message) error
    Subscribe(msgType string, handler HandlerFunc) error
    Unsubscribe(msgType string, handler HandlerFunc) error
}

type CommandBus interface {
    Send(ctx context.Context, cmd message.Command) error
}

type QueryBus interface {
    RequestReply(ctx context.Context, query message.Query) (message.Message, error)
}