// pkg/messagebus/middlewares.go

package messagebus

import (
    "context"
    "log"

	"github.com/mateusmacedo/golang-aether-sdk/core/pkg/application/message"
)

// LoggingMiddleware registra informações sobre as mensagens processadas.
func LoggingMiddleware(next HandlerFunc) HandlerFunc {
    return func(ctx context.Context, msg message.Message) error {
        log.Printf("Received message: %s", msg.Type())
        err := next(ctx, msg)
        if err != nil {
            log.Printf("Error processing message: %s", err)
        } else {
            log.Printf("Processed message: %s", msg.Type())
        }
        return err
    }
}

// AuthenticationMiddleware verifica a autenticidade da mensagem.
func AuthenticationMiddleware(next HandlerFunc) HandlerFunc {
    return func(ctx context.Context, msg message.Message) error {
        // Implementar lógica de autenticação aqui.
        // Se a autenticação falhar, retorne um erro.
        // Exemplo: if !authenticate(msg) { return errors.New("authentication failed") }

        return next(ctx, msg)
    }
}

// ValidationMiddleware verifica se a mensagem atende aos critérios de validação.
func ValidationMiddleware(next HandlerFunc) HandlerFunc {
    return func(ctx context.Context, msg message.Message) error {
        // Implementar lógica de validação aqui.
        // Se a validação falhar, retorne um erro.
        // Exemplo: if !validate(msg) { return errors.New("validation failed") }

        return next(ctx, msg)
    }
}