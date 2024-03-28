package messagebus

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/mateusmacedo/golang-aether-sdk/core/pkg/application/message"
)

type subscriptionInfo struct {
	handlers []HandlerFunc
	mutex    sync.RWMutex
}

type GenericMessageBus struct {
	subscriptions map[string]*subscriptionInfo
	middlewares   []Middleware
	mu            sync.RWMutex
}

func newSubscriptionInfo() *subscriptionInfo {
	return &subscriptionInfo{}
}

func (bus *GenericMessageBus) Use(middleware Middleware) {
	bus.middlewares = append(bus.middlewares, middleware)
}

func NewGenericMessageBus() *GenericMessageBus {
	return &GenericMessageBus{
		subscriptions: make(map[string]*subscriptionInfo),
	}
}

func (bus *GenericMessageBus) Publish(ctx context.Context, msg message.Message) error {
	bus.mu.RLock()
	subscription, ok := bus.subscriptions[msg.Type()]
	bus.mu.RUnlock()

	if !ok || len(subscription.handlers) == 0 {
		return fmt.Errorf("message type '%s' not handled", msg.Type())
	}

	// Aplicação dos middlewares e execução dos handlers em sequência.
	handler := func(ctxH context.Context, msgH message.Message) error {
		subscription.mutex.RLock()
		defer subscription.mutex.RUnlock()
		for _, h := range subscription.handlers {
			if err := h(ctxH, msgH); err != nil {
				log.Printf("Error executing handler for message type %s: %v", msgH.Type(), err)
				// TODO: Tratamento de erro e avaliação da necessidade de republicação.
			}
		}
		// TODO: Tratamento para acumulação de erros e retorno de erro único.
		return nil
	}

	// Encadeamento de middlewares.
	for i := len(bus.middlewares) - 1; i >= 0; i-- {
		handler = bus.middlewares[i](handler)
	}

	return handler(ctx, msg)
}

func (bus *GenericMessageBus) Subscribe(msgType string, handler HandlerFunc) error {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	if subscription, ok := bus.subscriptions[msgType]; !ok {
		bus.subscriptions[msgType] = newSubscriptionInfo()
		bus.subscriptions[msgType].handlers = append(bus.subscriptions[msgType].handlers, handler)
	} else {
		subscription.mutex.Lock()
		defer subscription.mutex.Unlock()
		subscription.handlers = append(subscription.handlers, handler)
	}

	return nil
}

func (bus *GenericMessageBus) Unsubscribe(msgType string, handler HandlerFunc) error {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	if subscription, ok := bus.subscriptions[msgType]; ok {
		subscription.mutex.Lock()
		defer subscription.mutex.Unlock()
		for i, h := range subscription.handlers {
			if &h == &handler {
				// TODO: Avaliar performance e necessidade de otimização para remoção de handler.
				subscription.handlers = append(subscription.handlers[:i], subscription.handlers[i+1:]...)
				return nil
			}
		}
	}

	return errors.New("handler not found for message type: " + msgType)
}

func (bus *GenericMessageBus) Send(ctx context.Context, cmd message.Command) error {
	bus.mu.RLock()
	subscription, ok := bus.subscriptions[cmd.Type()]
	bus.mu.RUnlock()

	if !ok || len(subscription.handlers) != 1 {
		return errors.New("exactly one handler expected for command type: " + cmd.Type())
	}

	return subscription.handlers[0](ctx, cmd)
}

func (bus *GenericMessageBus) RequestReply(ctx context.Context, query message.Query) (message.Message, error) {
	replyChan := make(chan message.Message, 1)

	// Definição de um handler temporário que espera pela resposta.
	// TODO: Avaliar o uso do pattern chain of responsibility para handlers temporários, e ou definir por injetor de dependências.
	handler := func(ctxH context.Context, msgH message.Message) error {
		select {
		case replyChan <- msgH:
			return nil
		case <-ctxH.Done():
			return ctxH.Err()
		}
	}

	// Registro do handler temporário.
	bus.Subscribe(query.ReplyType().Type(), handler)
	defer bus.Unsubscribe(query.ReplyType().Type(), handler) // Garante a remoção do handler temporário.

	err := bus.Publish(ctx, query)
	if err != nil {
		return nil, err
	}

	// Aguarda pela resposta dentro do prazo especificado.
	select {
	case reply := <-replyChan:
		return reply, nil
	case <-ctx.Done():
		switch ctx.Err() {
		case context.DeadlineExceeded:
			return nil, errors.New("timeout waiting for query reply")
		case context.Canceled:
			return nil, errors.New("request cancelled")
		}
		return nil, ctx.Err()
	}
}

// Garantia de conformidade com as interfaces através de verificações em tempo de compilação.
var _ MessageBus = &GenericMessageBus{}
var _ CommandBus = &GenericMessageBus{}
var _ QueryBus = &GenericMessageBus{}
