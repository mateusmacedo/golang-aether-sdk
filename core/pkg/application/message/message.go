package message

// Message é a interface base para todos os tipos de mensagens.
type Message interface {
    Type() string
}

// Command representa uma instrução para realizar uma ação.
type Command struct {
    // Campos específicos do comando
}

func (c Command) Type() string {
    return "command"
}

type Reply struct {
}

func (r Reply) Type() string {
    return "reply"
}

// Query representa uma solicitação de informação.
type Query struct {
    // Campos específicos da consulta
}

func (q Query) Type() string {
    return "query"
}

func (q Query) ReplyType() Message {
    return Reply{}
}

// Event representa um fato que ocorreu no sistema.
type Event struct {
    // Campos específicos do evento
}

func (e Event) Type() string {
    return "event"
}

func NewExampleCommand() Command {
    return Command{}
}

func NewExampleQuery() Query {
    return Query{}
}

func NewExampleReply() Reply {
    return Reply{}
}

func NewExampleEvent() Event {
    return Event{}
}