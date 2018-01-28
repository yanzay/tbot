package model

type MessageType int

const (
	MessageText MessageType = iota
	MessageContact
	MessageLocation
	MessageDocument
	MessageSticker
	MessagePhoto
	MessageAudio
	MessageKeyboard
	MessageSpecialKeyboard
)

const (
	ChatTypePrivate    = "private"
	ChatTypeGroup      = "group"
	ChatTypeSuperGroup = "supergroup"
	ChatTypeChannel    = "channel"
)

type Message struct {
	Type            MessageType
	Data            string
	Caption         string
	Replies         chan *Message
	From            User
	Contact         Contact
	Location        Location
	ChatID          int64
	ChatType        string
	DisablePreview  bool
	Markdown        bool
	Buttons         [][]string
	SpecialButtons  [][]KeyboardButton
	OneTimeKeyboard bool
	ForwardDate     int
}

// KeyboardButton is a button within a custom keyboard.
type KeyboardButton struct {
	Text            string `json:"text"`
	RequestContact  bool   `json:"request_contact"`
	RequestLocation bool   `json:"request_location"`
}

// Contact contains information about a contact.
//
// Note that LastName and UserID may be empty.
type Contact struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"` // optional
	UserID      int    `json:"user_id"`   // optional
}

// Location contains information about a place.
type Location struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

// User is a user on Telegram.
type User struct {
	ID           int    `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`     // optional
	UserName     string `json:"username"`      // optional
	LanguageCode string `json:"language_code"` // optional
}
