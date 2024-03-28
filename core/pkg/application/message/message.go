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

// Query representa uma solicitação de informação.
type Query struct {
    // Campos específicos da consulta
}

func (q Query) Type() string {
    return "query"
}

type QueryReply struct {
    // Campos específicos da resposta da consulta
}

func (qr QueryReply) Type() string {
    return "query_reply"
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

func NewExampleQueryReply() QueryReply {
    return QueryReply{}
}

func NewExampleEvent() Event {
    return Event{}
}