package messagebus_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/mateusmacedo/golang-aether-sdk/core/pkg/application/message"
	"github.com/mateusmacedo/golang-aether-sdk/core/pkg/application/messagebus"
)

func TestMessageBusPublish(t *testing.T) {
    bus := messagebus.NewGenericMessageBus()

    // Variável para verificar se o handler foi chamado.
    var handlerCalled bool
    var mu sync.Mutex

    // Definição do handler de teste para um comando específico.
    testHandler := func(ctx context.Context, msg message.Message) error {
        if _, ok := msg.(message.Command); !ok {
            t.Errorf("Message type is incorrect")
            return nil
        }

        mu.Lock()
        handlerCalled = true
        mu.Unlock()
        return nil
    }

    testCommand := message.NewExampleCommand()

    // Registro do handler de teste.
    err := bus.Subscribe(testCommand.Type(), testHandler)
    if err != nil {
        t.Fatalf("Subscribe failed: %v", err)
    }

    // Utiliza-se a instância específica da mensagem (Command, Query, Event) ao invés de uma mensagem genérica.
    err = bus.Publish(context.Background(), testCommand)
    if err != nil {
        t.Fatalf("Publish failed: %v", err)
    }

    // Espera para dar tempo ao handler de ser chamado.
    // Este atraso simula o tempo de processamento assíncrono.
    time.Sleep(100 * time.Millisecond)

    mu.Lock()
    defer mu.Unlock()
    if !handlerCalled {
        t.Errorf("handler was not called")
    }
}

func TestMessageBusRequestReply(t *testing.T) {
    bus := messagebus.NewGenericMessageBus()

    testQuery := message.NewExampleQuery()

    // Define um handler para a consulta que enviará a resposta esperada.
    bus.Subscribe(testQuery.Type(), func(ctx context.Context, msg message.Message) error {
        return bus.Publish(ctx, testQuery.ReplyType())
    })
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    reply, err := bus.RequestReply(ctx, testQuery)
    if err != nil {
        t.Fatalf("RequestReply failed: %v", err)
    }

    if reply.Type() != testQuery.ReplyType().Type() {
        t.Errorf("Expected reply type %v, got %v", testQuery.ReplyType().Type(), reply.Type())
    }
}

func TestMessageBusRequestReplyTimeout(t *testing.T) {
    bus := messagebus.NewGenericMessageBus()

    testQuery := message.NewExampleQuery()

    bus.Subscribe(testQuery.Type(), func(ctx context.Context, msg message.Message) error {
        time.Sleep(3* time.Millisecond)
        return nil
    })
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
    defer cancel()
    _, err := bus.RequestReply(ctx, testQuery)
    if err == nil {
        t.Fatalf("RequestReply did not return an error")
    }

    if err != nil && err.Error() != "timeout waiting for query reply" {
        t.Fatalf("Expected timeout error, got %v", err)
    }
}