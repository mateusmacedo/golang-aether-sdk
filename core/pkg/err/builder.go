package err

// ErrorBuilder é um construtor para erros customizados.
type ErrorBuilder struct {
	err *baseError
}

// New cria um novo ErrorBuilder com uma mensagem de erro.
func New(msg string) *ErrorBuilder {
	return &ErrorBuilder{
		err: &baseError{what: msg},
	}
}

// Wrap anexa uma causa ao erro.
func (b *ErrorBuilder) Wrap(cause error) *ErrorBuilder {
	b.err.cause = cause
	return b
}

// Build finaliza a construção do erro customizado.
func (b *ErrorBuilder) Build() CustomError {
	return b.err
}
