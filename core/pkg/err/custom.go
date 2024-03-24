package err

import (
	"fmt"
	"runtime"
	"time"
)

type CustomError interface {
	Error() string
	Unwrap() error
	WithStack([]string) CustomError
}

// baseError é uma implementação concreta de CustomError.
type baseError struct {
	what  string
	when  time.Time
	cause error
	stack []string
}

func (e *baseError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %v", e.what, e.cause)
	}
	return e.what
}

func (e *baseError) Unwrap() error {
	return e.cause
}

// WithStack captura a stack trace do ponto onde o erro foi criado ou decorado.
func (e *baseError) WithStack(s []string) CustomError {
	if e.stack == nil {
		e.stack = s
	}
	return e
}

// Skip captureStack and WithStack frames.
const SkipStackCaptureCalls = 2

func captureStack() []string {
	var stack []string
	for i := SkipStackCaptureCalls; true; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		stack = append(stack, fmt.Sprintf("%s:%d", file, line))
	}
	return stack
}
