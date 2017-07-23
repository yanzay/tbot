package model

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
	Type            MessageType
	Data            string
	Caption         string
	Replies         chan *Message
	From            User
	ChatID          int64
	DisablePreview  bool
	Markdown        bool
	Buttons         [][]string
	OneTimeKeyboard bool
	ForwardDate     int
}

// User is a user on Telegram.
type User struct {
	ID           int    `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`     // optional
	UserName     string `json:"username"`      // optional
	LanguageCode string `json:"language_code"` // optional
}
