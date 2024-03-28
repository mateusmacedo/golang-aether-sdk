package messagebus

import (
	"context"
	"time"

	"github.com/mateusmacedo/golang-aether-sdk/core/pkg/application/message"
)

// Middleware define a assinatura para um middleware no MessageBus.
type Middleware func(next HandlerFunc) HandlerFunc

// HandlerFunc define a assinatura para um manipulador de mensagens.
type HandlerFunc func(ctx context.Context, msg message.Message) error

// MessageBus define a interface para o barramento de mensagens com métodos básicos.
type MessageBus interface {
    Publish(ctx context.Context, msg message.Message) error
    Subscribe(msgType string, handler HandlerFunc) error
    Unsubscribe(msgType string, handler HandlerFunc) error
}

// CommandBus define a interface para o envio de comandos.
type CommandBus interface {
    Send(ctx context.Context, cmd message.Command) error
}

// QueryBus define a interface para requisições e respostas.
type QueryBus interface {
    RequestReply(ctx context.Context, query message.Query, timeout time.Duration) (message.Message, error)
}