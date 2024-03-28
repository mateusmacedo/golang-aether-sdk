package messagebus_test

import (
	"context"
	"testing"
	"time"

	"github.com/mateusmacedo/golang-aether-sdk/core/pkg/application/message"
	"github.com/mateusmacedo/golang-aether-sdk/core/pkg/application/messagebus"
)

func TestMessageBusPublish(t *testing.T) {
    t.Helper()
    var handlerCalled bool
    testCases := []struct {
        name          string
        message       message.Message
        handler       messagebus.HandlerFunc
        expectError   bool
        handlerCalled bool
    }{
        {
            name: "Test with ExampleCommand",
            message: message.NewExampleCommand(),
            handler: func(ctx context.Context, msg message.Message) error {
                if _, ok := msg.(message.Command); !ok {
                    t.Errorf("Message type is incorrect")
                    return nil
                }

                handlerCalled = true
                return nil
            },
            expectError:   false,
            handlerCalled: true,
        },
        {
            name: "Test with ExampleEvent",
            message: message.NewExampleEvent(),
            handler: func(ctx context.Context, msg message.Message) error {
                if _, ok := msg.(message.Event); !ok {
                    t.Errorf("Message type is incorrect")
                    return nil
                }

                handlerCalled = true
                return nil
            },
            expectError:   false,
            handlerCalled: true,
        },
        {
            name: "Test with ExampleQuery",
            message: message.NewExampleQuery(),
            handler: func(ctx context.Context, msg message.Message) error {
                if _, ok := msg.(message.Query); !ok {
                    t.Errorf("Message type is incorrect")
                    return nil
                }

                handlerCalled = true
                return nil
            },
            expectError:   false,
            handlerCalled: true,
        },
        // Adicione mais casos de teste conforme necessário
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(tr *testing.T) {
            handlerCalled = false
            bus := messagebus.NewGenericMessageBus()

            err := bus.Subscribe(tc.message.Type(), tc.handler)
            if err != nil {
                tr.Fatalf("Subscribe failed: %v", err)
            }

            err = bus.Publish(context.Background(), tc.message)
            if (err != nil) != tc.expectError {
                tr.Fatalf("Publish failed: %v", err)
            }

            time.Sleep(100 * time.Millisecond)

            if handlerCalled != tc.handlerCalled {
                tr.Errorf("handler was not called")
            }
        })
    }
}

func TestMessageBusRequestReply(t *testing.T) {
    t.Helper()
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
    t.Helper()
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