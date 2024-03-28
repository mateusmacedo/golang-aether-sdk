// main.go or a similar setup file

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/mateusmacedo/golang-aether-sdk/core/pkg/application/message"
	"github.com/mateusmacedo/golang-aether-sdk/core/pkg/application/messagebus"
)

func main() {
    // Criação do GenericMessageBus.
    bus := messagebus.NewGenericMessageBus()

    // Configuração dos middlewares (se necessário).

    // Exemplo de manipulador para um tipo de mensagem específico.
    bus.Subscribe("command", func(ctx context.Context, msg message.Message) error {
        // Lógica de manipulação da mensagem.
        fmt.Println("Handling Command:", msg)
        return nil
    })

    // Exemplo de publicação de uma mensagem.
    exampleMessage := message.NewExampleEvent() // Supondo que esta função crie uma nova mensagem de exemplo.
    if err := bus.Publish(context.Background(), exampleMessage); err != nil {
        fmt.Println("Error publishing message:", err)
    }

    // Exemplo de envio de um comando.
    exampleCommand := message.NewExampleCommand() // Supondo que esta função crie um novo comando de exemplo.
    if err := bus.Send(context.Background(), exampleCommand); err != nil {
        fmt.Println("Error sending command:", err)
    }

    // Exemplo de requisição e resposta.
    exampleQuery := message.NewExampleQuery() // Supondo que esta função crie uma nova consulta de exemplo.
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    reply, err := bus.RequestReply(ctx, exampleQuery)
    if err != nil {
        fmt.Println("Error requesting reply:", err)
    } else {
        fmt.Println("Received reply:", reply)
    }
}