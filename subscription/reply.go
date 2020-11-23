package subscription

type Reply = string

const (
	ReplyWelcome       Reply = "😸 Meow! Welcome to cat facts. You can unsubscribe at any time by replying stop."
	ReplyGoodbye             = "😿 Farewell from cat facts! You can re-subscribe at any time by replying start."
	ReplyHelp                = "🐾 This is cat facts. Reply start to start receiving a cat fact by text once daily. Reply stop to cancel your subscription."
	ReplyUnknown             = "🙀 Server cat could not understand your request! Reply help for more information."
	ReplyAlreadyExists       = "🙀 A subscription was already found for this number!"
)
