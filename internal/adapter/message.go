package adapter

type MessageType int

const (
	MessageText MessageType = iota
	MessageDocument
	MessageSticker
	MessagePhoto
	MessageAudio
	MessageKeyboard
)

type Message struct {
	Type           MessageType
	Data           string
	Caption        string
	Replies        chan *Message
	From           string
	ChatID         int64
	DisablePreview bool
	Buttons        [][]string
}
