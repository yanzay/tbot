package tbot

import (
	"encoding/json"
	"fmt"
)

// User is telegram user
type User struct {
	ID           int    `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
}

// ChatPhoto represents a chat photo
type ChatPhoto struct {
	SmallFileID string `json:"small_file_id"`
	BigFileID   string `json:"big_file_id"`
}

// Chat represents a chat
type Chat struct {
	ID                          string
	Type                        string
	Title                       string
	Username                    string
	FirstName                   string
	LastName                    string
	Photo                       *ChatPhoto
	Description                 string
	InviteLink                  string
	PinnedMessage               *Message
	StickerSetName              string
	AllMembersAreAdministrators bool
	CanSetStickerSet            bool
}

// UnmarshalJSON implements json.Unmarshaler
func (c *Chat) UnmarshalJSON(data []byte) error {
	s := &struct {
		ID                          int        `json:"id"`
		Type                        string     `json:"type"`
		Title                       string     `json:"title"`
		Username                    string     `json:"username"`
		FirstName                   string     `json:"first_name"`
		LastName                    string     `json:"last_name"`
		Photo                       *ChatPhoto `json:"photo"`
		Description                 string     `json:"description"`
		InviteLink                  string     `json:"invite_link"`
		PinnedMessage               *Message   `json:"pinned_message"`
		StickerSetName              string     `json:"sticker_set_name"`
		AllMembersAreAdministrators bool       `json:"all_members_are_administrators"`
		CanSetStickerSet            bool       `json:"can_set_sticker_set"`
	}{}
	err := json.Unmarshal(data, s)
	if err != nil {
		return err
	}
	*c = Chat{
		ID:                          fmt.Sprint(s.ID),
		Type:                        s.Type,
		Title:                       s.Title,
		Username:                    s.Username,
		FirstName:                   s.FirstName,
		LastName:                    s.LastName,
		Photo:                       s.Photo,
		Description:                 s.Description,
		InviteLink:                  s.InviteLink,
		PinnedMessage:               s.PinnedMessage,
		StickerSetName:              s.StickerSetName,
		AllMembersAreAdministrators: s.AllMembersAreAdministrators,
		CanSetStickerSet:            s.CanSetStickerSet,
	}
	return nil
}

// MessageEntity represents one special entity in a text message.
// For example, hashtags, usernames, URLs, etc.
type MessageEntity struct {
	Type   string `json:"type"`
	Offset int    `json:"offset"`
	Length int    `json:"length"`
	URL    string `json:"url"`
	User   *User  `json:"user"`
}

// Audio represents an audio file to be treated as music by the Telegram clients
type Audio struct {
	FileID    string `json:"file_id"`
	Duration  int    `json:"duration"`
	Performer string `json:"performer"`
	Title     string `json:"title"`
	MIMEType  string `json:"mime_type"`
	FileSize  int    `json:"file_size"`
}

// PhotoSize represents one size of a photo or a file/sticker thumbnail.
type PhotoSize struct {
	FileID   string `json:"file_id"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	FileSize int    `json:"file_size"`
}

// Document represents a general file
// (as opposed to photos, voice messages and audio files)
type Document struct {
	FileID   string     `json:"file_id"`
	Thumb    *PhotoSize `json:"thumb"`
	FileName string     `json:"file_name"`
	MIMEType string     `json:"mime_type"`
	FileSize int        `json:"file_size"`
}

// Game represents a game. Use BotFather to create and edit games,
// their short names will act as unique identifiers.
type Game struct {
	Title        string           `json:"title"`
	Description  string           `json:"description"`
	Photo        []*PhotoSize     `json:"photo"`
	Text         string           `json:"text"`
	TextEntities []*MessageEntity `json:"text_entities"`
	Animation    *Animation       `json:"animation"`
}

// Animation represents an animation file
// to be displayed in the message containing a game
type Animation struct {
	FileID   string     `json:"file_id"`
	Thumb    *PhotoSize `json:"thumb"`
	FileName string     `json:"file_name"`
	MimeType string     `json:"mime_type"`
	FileSize int        `json:"file_size"`
}

// Sticker represents a sticker
type Sticker struct {
	FileID       string        `json:"file_id"`
	Width        int           `json:"width"`
	Height       int           `json:"height"`
	IsAnimated   bool          `json:"is_animated"`
	Thumb        *PhotoSize    `json:"thumb"`
	Emoji        string        `json:"emoji"`
	MaskPosition *MaskPosition `json:"mask_position"`
	SetName      string        `json:"set_name"`
	FileSize     int           `json:"file_size"`
}

// MaskPosition describes the position on faces
// where a mask should be placed by default
type MaskPosition struct {
	Point  string  `json:"point"`
	XShift float32 `json:"x_shift"`
	YShift float32 `json:"y_shift"`
	Scale  float32 `json:"scale"`
}

// Video represents a video file
type Video struct {
	FileID    string     `json:"file_id"`
	Width     int        `json:"width"`
	Height    int        `json:"height"`
	Duration  int        `json:"duration"`
	Thumbnail *PhotoSize `json:"thumb"`
	MimeType  string     `json:"mime_type"`
	FileSize  int        `json:"file_size"`
}

// Voice represents a voice note
type Voice struct {
	FileID   string `json:"file_id"`
	Duration int    `json:"duration"`
	MimeType string `json:"mime_type"`
	FileSize int    `json:"file_size"`
}

// VideoNote represents a video message
type VideoNote struct {
	FileID   string     `json:"file_id"`
	Length   int        `json:"length"`
	Duration int        `json:"duration"`
	Thumb    *PhotoSize `json:"thumb"`
	FileSize int        `json:"file_size"`
}

// Contact represents a phone contact
type Contact struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	UserID      int    `json:"user_id"`
}

// Location represents a point on the map
type Location struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

// Venue represents a venue
type Venue struct {
	Location     Location `json:"location"`
	Title        string   `json:"title"`
	Address      string   `json:"address"`
	FoursquareID string   `json:"foursquare_id"`
}

// Invoice contains basic information about an invoice
type Invoice struct {
	Title          string `json:"title"`
	Description    string `json:"description"`
	StartParameter string `json:"start_parameter"`
	Currency       string `json:"currency"`
	TotalAmount    int    `json:"total_amount"`
}

// SuccessfulPayment contains basic information about a successful payment
type SuccessfulPayment struct {
	Currency                string     `json:"currency"`
	TotalAmount             int        `json:"total_amount"`
	InvoicePayload          string     `json:"invoice_payload"`
	ShippingOptionID        string     `json:"shipping_option_id"`
	OrderInfo               *OrderInfo `json:"order_info"`
	TelegramPaymentChargeID string     `json:"telegram_payment_charge_id"`
	ProviderPaymentChargeID string     `json:"provider_payment_charge_id"`
}

// OrderInfo represents information about an order
type OrderInfo struct {
	Name            string           `json:"name"`
	PhoneNumber     string           `json:"phone_number"`
	Email           string           `json:"email"`
	ShippingAddress *ShippingAddress `json:"shipping_address"`
}

// ShippingAddress represents a shipping address
type ShippingAddress struct {
	CountryCode string `json:"country_code"`
	State       string `json:"state"`
	City        string `json:"city"`
	StreetLine1 string `json:"street_line1"`
	StreetLine2 string `json:"street_line2"`
	PostCode    string `json:"post_code"`
}

// Message represents a message
type Message struct {
	MessageID             int                   `json:"message_id"`
	From                  *User                 `json:"from"`
	Date                  int64                 `json:"date"`
	Chat                  Chat                  `json:"chat"`
	ForwardFrom           *User                 `json:"forward_from"`
	ForwardFromChat       *Chat                 `json:"forward_from_chat"`
	ForwardFromMessageID  int                   `json:"forward_from_message_id"`
	ForwardSignature      string                `json:"forward_signature"`
	ForwardSenderName     string                `json:"forward_sender_name"`
	ForwardDate           int64                 `json:"forward_date"`
	ReplyToMessage        *Message              `json:"reply_to_message"`
	EditDate              int64                 `json:"edit_date"`
	MediaGroupID          string                `json:"media_group_id"`
	AuthorSignature       string                `json:"author_signature"`
	Text                  string                `json:"text"`
	Entities              []*MessageEntity      `json:"entities"`
	CaptionEntities       []*MessageEntity      `json:"caption_entities"`
	Audio                 *Audio                `json:"audio"`
	Document              *Document             `json:"document"`
	Game                  *Game                 `json:"game"`
	Photo                 []*PhotoSize          `json:"photo"`
	Sticker               *Sticker              `json:"sticker"`
	Video                 *Video                `json:"video"`
	Voice                 *Voice                `json:"voice"`
	VideoNote             *VideoNote            `json:"video_note"`
	Caption               string                `json:"caption"`
	Contact               *Contact              `json:"contact"`
	Location              *Location             `json:"location"`
	Venue                 *Venue                `json:"venue"`
	Poll                  *Poll                 `json:"poll"`
	NewChatMembers        []*User               `json:"new_chat_members"`
	LeftChatMember        *User                 `json:"left_chat_member"`
	NewChatTitle          string                `json:"new_chat_title"`
	NewChatPhoto          []*PhotoSize          `json:"new_chat_photo"`
	DeleteChatPhoto       bool                  `json:"delete_chat_photo"`
	GroupChatCreated      bool                  `json:"group_chat_created"`
	SupergroupChatCreated bool                  `json:"supergroup_chat_created"`
	ChannelChatCreated    bool                  `json:"channel_chat_created"`
	MigrateToChatID       int                   `json:"migrate_to_chat_id"`
	MigrateFromChatID     int                   `json:"migrate_from_chat_id"`
	PinnedMessage         *Message              `json:"pinned_message"`
	Invoice               *Invoice              `json:"invoice"`
	SuccessfulPayment     *SuccessfulPayment    `json:"successful_payment"`
	ConnectedWebsite      string                `json:"connected_website"`
	PassportData          *PassportData         `json:"passport_data"`
	ReplyMarkup           *InlineKeyboardMarkup `json:"reply_markup"`
}

// InlineQuery represents an incoming inline query
type InlineQuery struct {
	ID       string    `json:"id"`
	From     *User     `json:"from"`
	Location *Location `json:"location"`
	Query    string    `json:"query"`
	Offset   string    `json:"offset"`
}

// ChosenInlineResult represents a result of an inline query
// that was chosen by the user and sent to their chat partner
type ChosenInlineResult struct {
	ResultID        string    `json:"result_id"`
	From            *User     `json:"from"`
	Location        *Location `json:"location"`
	InlineMessageID string    `json:"inline_message_id"`
	Query           string    `json:"query"`
}

// CallbackQuery represents an incoming callback query
// from a callback button in an inline keyboard
type CallbackQuery struct {
	ID              string   `json:"id"`
	From            *User    `json:"from"`
	Message         *Message `json:"message"`
	InlineMessageID string   `json:"inline_message_id"`
	ChatInstance    string   `json:"chat_instance"`
	Data            string   `json:"data"`
	GameShortName   string   `json:"game_short_name"`
}

// ShippingQuery contains information about an incoming shipping query
type ShippingQuery struct {
	ID              string           `json:"id"`
	From            *User            `json:"from"`
	InvoicePayload  string           `json:"invoice_payload"`
	ShippingAddress *ShippingAddress `json:"shipping_address"`
}

// PreCheckoutQuery contains information about an incoming pre-checkout query
type PreCheckoutQuery struct {
	ID               string     `json:"id"`
	From             *User      `json:"from"`
	Currency         string     `json:"currency"`
	TotalAmount      int        `json:"total_amount"`
	InvoicePayload   string     `json:"invoice_payload"`
	ShippingOptionID string     `json:"shipping_option_id"`
	OrderInfo        *OrderInfo `json:"order_info"`
}

// Update represents an incoming update
// UpdateID is unique identifier
// At most one of the other fields can be not nil
type Update struct {
	UpdateID           int                 `json:"update_id"`
	Message            *Message            `json:"message"`
	EditedMessage      *Message            `json:"edited_message"`
	ChannelPost        *Message            `json:"channel_post"`
	EditedChannelPost  *Message            `json:"edited_channel_post"`
	InlineQuery        *InlineQuery        `json:"inline_query"`
	ChosenInlineResult *ChosenInlineResult `json:"chosen_inline_result"`
	CallbackQuery      *CallbackQuery      `json:"callback_query"`
	ShippingQuery      *ShippingQuery      `json:"shipping_query"`
	PreCheckoutQuery   *PreCheckoutQuery   `json:"pre_checkout_query"`
	Poll               *Poll               `json:"poll"`
}

// PassportData contains information about Telegram Passport data shared with the bot by the user
type PassportData struct {
	Data        []EncryptedPassportElement `json:"data"`
	Credentials EncryptedCredentials       `json:"credentials"`
}

// EncryptedPassportElement contains information about documents or other Telegram Passport elements shared with the bot by the user
type EncryptedPassportElement struct {
	Type        string          `json:"type"`
	Data        string          `json:"data"`
	PhoneNumber string          `json:"phone_number"`
	Email       string          `json:"email"`
	Files       []*PassportFile `json:"files"`
	FrontSide   *PassportFile   `json:"front_side"`
	ReverseSide *PassportFile   `json:"reverse_side"`
	Selfie      *PassportFile   `json:"selfie"`
}

// PassportFile represents a file uploaded to Telegram Passport
type PassportFile struct {
	FileID   string `json:"file_id"`
	FileSize int    `json:"file_size"`
	FileDate int    `json:"file_date"`
}

// EncryptedCredentials contains data required for decrypting and authenticating EncryptedPassportElement
type EncryptedCredentials struct {
	Data   string `json:"data"`
	Hash   string `json:"hash"`
	Secret string `json:"secret"`
}

// Poll represents native telegram poll
type Poll struct {
	ID       string       `json:"id"`
	Question string       `json:"question"`
	Options  []PollOption `json:"options"`
	IsClosed bool         `json:"is_closed"`
}

// PollOption is an option for Poll
type PollOption struct {
	Text       string `json:"text"`
	VoterCount int    `json:"voter_count"`
}
