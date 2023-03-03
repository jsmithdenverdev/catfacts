package catfacts

type Command string

const (
	CommandSubscribe   Command = "subscribe"
	CommandUnsubscribe Command = "unsubscribe"
	CommandCuss        Command = "cuss"
)
