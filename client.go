package tbot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Client is a low-level Telegram client
type Client struct {
	token         string
	url           string
	httpClient    *http.Client
	nextOffset    int
	logger        Logger
	bufferSize    int
	timeout       int
	updatesParams url.Values
}

// NewClient creates new Telegram API client
func NewClient(token string, httpClient *http.Client, baseURL string) *Client {
	return &Client{
		token:      token,
		httpClient: httpClient,
		url:        fmt.Sprintf("%s/bot%s/", baseURL, token) + "%s",
	}
}

type inputFile struct {
	field string
	name  string
}

type sendOption func(url.Values)

// Generic message options
var (
	OptParseModeHTML = func(r url.Values) {
		r.Set("parse_mode", "HTML")
	}
	OptParseModeMarkdown = func(r url.Values) {
		r.Set("parse_mode", "Markdown")
	}
	OptDisableNotification = func(r url.Values) {
		r.Set("disable_notification", "true")
	}
	OptReplyToMessageID = func(id int) sendOption {
		return func(r url.Values) {
			r.Set("reply_to_message_id", strconv.Itoa(id))
		}
	}
)

func structString(s interface{}) string {
	str, _ := json.Marshal(s)
	return string(str)
}

// GetMe returns info about bot as a User object
func (c *Client) GetMe() (*User, error) {
	me := &User{}
	err := c.doRequest("getMe", nil, me)
	return me, err
}

type forceReply struct {
	ForceReply bool `json:"force_reply"`
	Selective  bool `json:"selective"`
}

type replyKeyboardRemove struct {
	RemoveKeyboard bool `json:"remove_keyboard"`
	Selective      bool `json:"selective"`
}

// InlineKeyboardMarkup represents an inline keyboard that appears right next to the message it belongs to
type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

// InlineKeyboardButton represents one button of an inline keyboard
type InlineKeyboardButton struct {
	Text                         string    `json:"text"`
	URL                          string    `json:"url,omitempty"`
	LoginURL                     *LoginURL `json:"login_url,omitempty"`
	CallbackData                 string    `json:"callback_data,omitempty"`
	SwitchInlineQuery            *string   `json:"switch_inline_query,omitempty"`
	SwitchInlineQueryCurrentChat *string   `json:"switch_inline_query_current_chat,omitempty"`
}

// LoginURL is a property of InlineKeyboardButton for Seamless Login feature
type LoginURL struct {
	URL                string  `json:"url"`
	ForwardText        *string `json:"forward_text,omitempty"`
	BotUsername        *string `json:"bot_username,omitempty"`
	RequestWriteAccess *string `json:"request_write_access,omitempty"`
}

// ReplyKeyboardMarkup represents a custom keyboard with reply options
type ReplyKeyboardMarkup struct {
	Keyboard        [][]KeyboardButton `json:"keyboard"`
	ResizeKeyboard  bool               `json:"resize_keyboard"`
	OneTimeKeyboard bool               `json:"one_time_keyboard"`
	Selective       bool               `json:"selective"`
}

// KeyboardButton represents one button of the reply keyboard
type KeyboardButton struct {
	Text            string `json:"text"`
	RequestContact  bool   `json:"request_contact"`
	RequestLocation bool   `json:"request_location"`
}

func (c *Client) setWebhook(webhookURL string) error {
	req := url.Values{}
	req.Set("url", webhookURL)
	var set bool
	return c.doRequest("setWebhook", req, &set)
}

func (c *Client) deleteWebhook() error {
	var ok bool
	return c.doRequest("deleteWebhook", url.Values{}, &ok)
}

// SendMessage options
var (
	OptDisableWebPagePreview = func(r url.Values) {
		r.Set("disable_web_page_preview", "true")
	}
	OptInlineKeyboardMarkup = func(markup *InlineKeyboardMarkup) sendOption {
		return func(r url.Values) {
			r.Set("reply_markup", structString(markup))
		}
	}
	OptReplyKeyboardMarkup = func(markup *ReplyKeyboardMarkup) sendOption {
		return func(r url.Values) {
			r.Set("reply_markup", structString(markup))
		}
	}
	OptReplyKeyboardRemove = func(r url.Values) {
		r.Set("reply_markup", structString(&replyKeyboardRemove{RemoveKeyboard: true}))
	}
	OptReplyKeyboardRemoveSelective = func(r url.Values) {
		r.Set("reply_markup", structString(&replyKeyboardRemove{RemoveKeyboard: true, Selective: true}))
	}
	OptForceReply = func(r url.Values) {
		r.Set("reply_markup", structString(&forceReply{ForceReply: true}))
	}
	OptForceReplySelective = func(r url.Values) {
		r.Set("reply_markup", structString(&forceReply{ForceReply: true, Selective: true}))
	}
)

/*
SendMessage sends message to telegram chat. Available options:
	- OptParseModeHTML
	- OptParseModeMarkdown
	- OptDisableWebPagePreview
	- OptDisableNotification
	- OptReplyToMessageID(id int)
	- OptInlineKeyboardMarkup(markup *InlineKeyboardMarkup)
	- OptReplyKeyboardMarkup(markup *ReplyKeyboardMarkup)
	- OptReplyKeyboardRemove
	- OptReplyKeyboardRemoveSelective
	- OptForceReply
	- OptForceReplySelective
*/
func (c *Client) SendMessage(chatID string, text string, opts ...sendOption) (*Message, error) {
	req := url.Values{}
	req.Set("chat_id", chatID)
	req.Set("text", text)
	for _, opt := range opts {
		opt(req)
	}
	msg := &Message{}
	err := c.doRequest("sendMessage", req, msg)
	return msg, err
}

/*
ForwardMessage forwards message from one chat to another. Available options:
	- OptDisableNotification
*/
func (c *Client) ForwardMessage(chatID, fromChatID string, messageID int, opts ...sendOption) (*Message, error) {
	req := url.Values{}
	req.Set("chat_id", chatID)
	req.Set("from_chat_id", fromChatID)
	req.Set("message_id", strconv.Itoa(messageID))
	for _, opt := range opts {
		opt(req)
	}
	msg := &Message{}
	err := c.doRequest("forwardMessage", req, msg)
	return msg, err
}

// SendAudio options
var (
	OptDuration = func(duration int) sendOption {
		return func(r url.Values) {
			r.Set("duration", strconv.Itoa(duration))
		}
	}
	OptPerformer = func(performer string) sendOption {
		return func(r url.Values) {
			r.Set("performer", performer)
		}
	}
	OptTitle = func(title string) sendOption {
		return func(r url.Values) {
			r.Set("title", title)
		}
	}
)

/*
SendAudio sends pre-uploaded audio to the chat. Pass fileID of the uploaded file. Available options:
	- OptCaption(caption string)
	- OptDuration(duration int)
	- OptPerformer(performer string)
	- OptTitle(title string)
	- OptParseModeHTML
	- OptParseModeMarkdown
	- OptDisableNotification
	- OptReplyToMessageID(id int)
	- OptInlineKeyboardMarkup(markup *InlineKeyboardMarkup)
	- OptReplyKeyboardMarkup(markup *ReplyKeyboardMarkup)
	- OptReplyKeyboardRemove
	- OptReplyKeyboardRemoveSelective
	- OptForceReply
	- OptForceReplySelective
*/
func (c *Client) SendAudio(chatID string, fileID string, opts ...sendOption) (*Message, error) {
	req := url.Values{}
	req.Set("chat_id", chatID)
	req.Set("audio", fileID)
	for _, opt := range opts {
		opt(req)
	}
	msg := &Message{}
	err := c.doRequest("sendAudio", req, msg)
	return msg, err
}

/*
SendAudioFile sends file contents as an audio to the chat. Pass filename to send. Available options:
	- OptCaption(caption string)
	- OptDuration(duration int)
	- OptPerformer(performer string)
	- OptTitle(title string)
	- OptParseModeHTML
	- OptParseModeMarkdown
	- OptDisableNotification
	- OptReplyToMessageID(id int)
	- OptInlineKeyboardMarkup(markup *InlineKeyboardMarkup)
	- OptReplyKeyboardMarkup(markup *ReplyKeyboardMarkup)
	- OptReplyKeyboardRemove
	- OptReplyKeyboardRemoveSelective
	- OptForceReply
	- OptForceReplySelective
*/
func (c *Client) SendAudioFile(chatID string, filename string, opts ...sendOption) (*Message, error) {
	req := url.Values{}
	req.Set("chat_id", chatID)
	for _, opt := range opts {
		opt(req)
	}
	msg := &Message{}
	err := c.doRequestWithFiles("sendAudio", req, msg, inputFile{field: "audio", name: filename})
	return msg, err
}

// SendPhoto options
var (
	OptCaption = func(caption string) sendOption {
		return func(r url.Values) {
			r.Set("caption", caption)
		}
	}
)

/*
SendPhoto sends pre-uploaded photo to the chat. Pass fileID of the photo. Available options:
	- OptCaption(caption string)
	- OptParseModeHTML
	- OptParseModeMarkdown
	- OptDisableNotification
	- OptReplyToMessageID(id int)
	- OptInlineKeyboardMarkup(markup *InlineKeyboardMarkup)
	- OptReplyKeyboardMarkup(markup *ReplyKeyboardMarkup)
	- OptReplyKeyboardRemove
	- OptReplyKeyboardRemoveSelective
	- OptForceReply
	- OptForceReplySelective
*/
func (c *Client) SendPhoto(chatID string, fileID string, opts ...sendOption) (*Message, error) {
	req := url.Values{}
	req.Set("chat_id", chatID)
	req.Set("photo", fileID)
	for _, opt := range opts {
		opt(req)
	}
	msg := &Message{}
	err := c.doRequest("sendPhoto", req, msg)
	return msg, err
}

/*
SendPhotoFile sends photo file contents to the chat. Pass filename to send. Available options:
	- OptCaption(caption string)
	- OptParseModeHTML
	- OptParseModeMarkdown
	- OptDisableNotification
	- OptReplyToMessageID(id int)
	- OptInlineKeyboardMarkup(markup *InlineKeyboardMarkup)
	- OptReplyKeyboardMarkup(markup *ReplyKeyboardMarkup)
	- OptReplyKeyboardRemove
	- OptReplyKeyboardRemoveSelective
	- OptForceReply
	- OptForceReplySelective
*/
func (c *Client) SendPhotoFile(chatID string, filename string, opts ...sendOption) (*Message, error) {
	req := url.Values{}
	req.Set("chat_id", chatID)
	for _, opt := range opts {
		opt(req)
	}
	msg := &Message{}
	err := c.doRequestWithFiles("sendPhoto", req, msg, inputFile{field: "photo", name: filename})
	return msg, err
}

/*
SendDocument sends document to the chat. Pass fileID of the document. Available options:
	- OptCaption(caption string)
	- OptParseModeHTML
	- OptParseModeMarkdown
	- OptDisableNotification
	- OptReplyToMessageID(id int)
	- OptInlineKeyboardMarkup(markup *InlineKeyboardMarkup)
	- OptReplyKeyboardMarkup(markup *ReplyKeyboardMarkup)
	- OptReplyKeyboardRemove
	- OptReplyKeyboardRemoveSelective
	- OptForceReply
	- OptForceReplySelective
*/
func (c *Client) SendDocument(chatID string, fileID string, opts ...sendOption) (*Message, error) {
	req := url.Values{}
	req.Set("chat_id", chatID)
	req.Set("document", fileID)
	for _, opt := range opts {
		opt(req)
	}
	msg := &Message{}
	err := c.doRequest("sendDocument", req, msg)
	return msg, err
}

/*
SendDocumentFile sends document file contents to the chat. Pass filename to send. Available options:
	- OptCaption(caption string)
	- OptParseModeHTML
	- OptParseModeMarkdown
	- OptDisableNotification
	- OptReplyToMessageID(id int)
	- OptInlineKeyboardMarkup(markup *InlineKeyboardMarkup)
	- OptReplyKeyboardMarkup(markup *ReplyKeyboardMarkup)
	- OptReplyKeyboardRemove
	- OptReplyKeyboardRemoveSelective
	- OptForceReply
	- OptForceReplySelective
*/
func (c *Client) SendDocumentFile(chatID string, filename string, opts ...sendOption) (*Message, error) {
	req := url.Values{}
	req.Set("chat_id", chatID)
	for _, opt := range opts {
		opt(req)
	}
	msg := &Message{}
	err := c.doRequestWithFiles("sendDocument", req, msg, inputFile{field: "document", name: filename})
	return msg, err
}

// SendVideo options
var (
	OptWidth = func(width int) sendOption {
		return func(r url.Values) {
			r.Set("width", strconv.Itoa(width))
		}
	}
	OptHeight = func(height int) sendOption {
		return func(r url.Values) {
			r.Set("height", strconv.Itoa(height))
		}
	}
	OptSupportsStreaming = func(r url.Values) {
		r.Set("supports_streaming", "true")
	}
)

/*
SendVideo sends pre-uploaded video to chat. Pass fileID of the uploaded video. Available options:
	- OptDuration(duration int)
	- OptWidth(width int)
	- OptHeight(height int)
	- OptSupportsStreaming
	- OptCaption(caption string)
	- OptParseModeHTML
	- OptParseModeMarkdown
	- OptDisableNotification
	- OptReplyToMessageID(id int)
	- OptInlineKeyboardMarkup(markup *InlineKeyboardMarkup)
	- OptReplyKeyboardMarkup(markup *ReplyKeyboardMarkup)
	- OptReplyKeyboardRemove
	- OptReplyKeyboardRemoveSelective
	- OptForceReply
	- OptForceReplySelective
*/
func (c *Client) SendVideo(chatID string, fileID string, opts ...sendOption) (*Message, error) {
	req := url.Values{}
	req.Set("chat_id", chatID)
	req.Set("video", fileID)
	for _, opt := range opts {
		opt(req)
	}
	msg := &Message{}
	err := c.doRequest("sendVideo", req, msg)
	return msg, err
}

/*
SendVideoFile sends video file contents to the chat. Pass filename to send. Available options:
	- OptDuration(duration int)
	- OptWidth(width int)
	- OptHeight(height int)
	- OptSupportsStreaming
	- OptCaption(caption string)
	- OptParseModeHTML
	- OptParseModeMarkdown
	- OptDisableNotification
	- OptReplyToMessageID(id int)
	- OptInlineKeyboardMarkup(markup *InlineKeyboardMarkup)
	- OptReplyKeyboardMarkup(markup *ReplyKeyboardMarkup)
	- OptReplyKeyboardRemove
	- OptReplyKeyboardRemoveSelective
	- OptForceReply
	- OptForceReplySelective
*/
func (c *Client) SendVideoFile(chatID string, filename string, opts ...sendOption) (*Message, error) {
	req := url.Values{}
	req.Set("chat_id", chatID)
	for _, opt := range opts {
		opt(req)
	}
	msg := &Message{}
	err := c.doRequestWithFiles("sendVideo", req, msg, inputFile{field: "video", name: filename})
	return msg, err
}

// SendAnimation options
var (
	OptThumb = func(filename string) sendOption {
		return func(v url.Values) {
			v.Set("thumb", filename)
		}
	}
)

/*
SendAnimation sends animation to chat. Pass fileID to send. Available options:
	- OptDuration(duration int)
	- OptWidth(width int)
	- OptHeight(height int)
	- OptThumb(filename string)
	- OptCaption(caption string)
	- OptParseModeHTML
	- OptParseModeMarkdown
	- OptDisableNotification
	- OptReplyToMessageID(id int)
	- OptInlineKeyboardMarkup(markup *InlineKeyboardMarkup)
	- OptReplyKeyboardMarkup(markup *ReplyKeyboardMarkup)
	- OptReplyKeyboardRemove
	- OptReplyKeyboardRemoveSelective
	- OptForceReply
	- OptForceReplySelective
*/
func (c *Client) SendAnimation(chatID string, fileID string, opts ...sendOption) (*Message, error) {
	req := url.Values{}
	req.Set("chat_id", chatID)
	req.Set("animation", fileID)
	for _, opt := range opts {
		opt(req)
	}
	msg := &Message{}
	var err error
	if len(req.Get("thumb")) > 0 {
		thumb := req.Get("thumb")
		req.Del("thumb")
		err = c.doRequestWithFiles("sendAnimation", req, msg, inputFile{field: "thumb", name: thumb})
	} else {
		err = c.doRequest("sendAnimation", req, msg)
	}
	return msg, err
}

/*
SendAnimationFile sends animation file contents to the chat. Pass filename to send. Available options:
	- OptDuration(duration int)
	- OptWidth(width int)
	- OptHeight(height int)
	- OptThumb(filename string)
	- OptCaption(caption string)
	- OptParseModeHTML
	- OptParseModeMarkdown
	- OptDisableNotification
	- OptReplyToMessageID(id int)
	- OptInlineKeyboardMarkup(markup *InlineKeyboardMarkup)
	- OptReplyKeyboardMarkup(markup *ReplyKeyboardMarkup)
	- OptReplyKeyboardRemove
	- OptReplyKeyboardRemoveSelective
	- OptForceReply
	- OptForceReplySelective
*/
func (c *Client) SendAnimationFile(chatID string, filename string, opts ...sendOption) (*Message, error) {
	req := url.Values{}
	req.Set("chat_id", chatID)
	for _, opt := range opts {
		opt(req)
	}
	msg := &Message{}
	files := []inputFile{{field: "animation", name: filename}}
	if len(req.Get("thumb")) > 0 {
		thumb := req.Get("thumb")
		req.Del("thumb")
		files = append(files, inputFile{field: "thumb", name: thumb})
	}
	err := c.doRequestWithFiles("sendAnimation", req, msg, files...)
	return msg, err
}

/*
SendVoice sends audio file as a voice message. Pass file_id of previously uploaded file. Available options:
	- OptCaption(caption string)
	- OptDuration(duration int)
	- OptParseModeHTML
	- OptParseModeMarkdown
	- OptDisableNotification
	- OptReplyToMessageID(id int)
	- OptInlineKeyboardMarkup(markup *InlineKeyboardMarkup)
	- OptReplyKeyboardMarkup(markup *ReplyKeyboardMarkup)
	- OptReplyKeyboardRemove
	- OptReplyKeyboardRemoveSelective
	- OptForceReply
	- OptForceReplySelective
*/
func (c *Client) SendVoice(chatID string, fileID string, opts ...sendOption) (*Message, error) {
	req := url.Values{}
	req.Set("chat_id", chatID)
	req.Set("voice", fileID)
	for _, opt := range opts {
		opt(req)
	}
	msg := &Message{}
	err := c.doRequest("sendVoice", req, msg)
	return msg, err
}

/*
SendVoiceFile sends the audio file as a voice message. Pass filename to send. Available options:
	- OptCaption(caption string)
	- OptDuration(duration int)
	- OptParseModeHTML
	- OptParseModeMarkdown
	- OptDisableNotification
	- OptReplyToMessageID(id int)
	- OptInlineKeyboardMarkup(markup *InlineKeyboardMarkup)
	- OptReplyKeyboardMarkup(markup *ReplyKeyboardMarkup)
	- OptReplyKeyboardRemove
	- OptReplyKeyboardRemoveSelective
	- OptForceReply
	- OptForceReplySelective
*/
func (c *Client) SendVoiceFile(chatID string, filename string, opts ...sendOption) (*Message, error) {
	req := url.Values{}
	req.Set("chat_id", chatID)
	for _, opt := range opts {
		opt(req)
	}
	msg := &Message{}
	err := c.doRequestWithFiles("sendVoice", req, msg, inputFile{field: "voice", name: filename})
	return msg, err
}

// SendVideoNote options
var (
	OptLength = func(length int) sendOption {
		return func(v url.Values) {
			v.Set("length", fmt.Sprint(length))
		}
	}
)

/*
SendVideoNote sends video note. Pass fileID of previously uploaded video note. Available options:
	- OptDuration(duration int)
	- OptLength(length int)
	- OptThumb(filename string)
	- OptDisableNotification
	- OptReplyToMessageID(id int)
	- OptInlineKeyboardMarkup(markup *InlineKeyboardMarkup)
	- OptReplyKeyboardMarkup(markup *ReplyKeyboardMarkup)
	- OptReplyKeyboardRemove
	- OptReplyKeyboardRemoveSelective
	- OptForceReply
	- OptForceReplySelective
*/
func (c *Client) SendVideoNote(chatID string, fileID string, opts ...sendOption) (*Message, error) {
	req := url.Values{}
	req.Set("chat_id", chatID)
	req.Set("video_note", fileID)
	for _, opt := range opts {
		opt(req)
	}
	msg := &Message{}
	var err error
	if len(req.Get("thumb")) > 0 {
		thumb := req.Get("thumb")
		req.Del("thumb")
		err = c.doRequestWithFiles("sendVideoNote", req, msg, inputFile{field: "thumb", name: thumb})
	} else {
		err = c.doRequest("sendVideoNote", req, msg)
	}
	return msg, err
}

/*
SendVideoNoteFile sends video note to chat. Pass filename to upload. Available options:
	- OptDuration(duration int)
	- OptLength(length int)
	- OptThumb(filename string)
	- OptDisableNotification
	- OptReplyToMessageID(id int)
	- OptInlineKeyboardMarkup(markup *InlineKeyboardMarkup)
	- OptReplyKeyboardMarkup(markup *ReplyKeyboardMarkup)
	- OptReplyKeyboardRemove
	- OptReplyKeyboardRemoveSelective
	- OptForceReply
	- OptForceReplySelective
*/
func (c *Client) SendVideoNoteFile(chatID string, filename string, opts ...sendOption) (*Message, error) {
	req := url.Values{}
	req.Set("chat_id", chatID)
	for _, opt := range opts {
		opt(req)
	}
	files := []inputFile{{field: "video_note", name: filename}}
	if len(req.Get("thumb")) > 0 {
		thumb := req.Get("thumb")
		req.Del("thumb")
		files = append(files, inputFile{field: "thumb", name: thumb})
	}
	msg := &Message{}
	err := c.doRequestWithFiles("sendVideoNote", req, msg, files...)
	return msg, err
}

// InputMedia file
type InputMedia interface {
	inputMedia()
}

var (
	_ InputMedia = InputMediaPhoto{}
	_ InputMedia = InputMediaVideo{}
)

// InputMediaPhoto represents a photo to be sent
type InputMediaPhoto struct {
	Type      string `json:"type"`
	Media     string `json:"media"`
	Caption   string `json:"caption,omitempty"`
	ParseMode string `json:"parse_mode,omitempty"`
}

func (InputMediaPhoto) inputMedia() {}

// InputMediaVideo represents a video to be sent
type InputMediaVideo struct {
	Type              string `json:"type"`
	Media             string `json:"media"`
	Thumb             string `json:"thumb,omitempty"`
	Caption           string `json:"caption,omitempty"`
	ParseMode         string `json:"parse_mode,omitempty"`
	Width             int    `json:"width,omitempty"`
	Height            int    `json:"height,omitempty"`
	Duration          int    `json:"duration,omitempty"`
	SupportsStreaming bool   `json:"supports_streaming,omitempty"`
}

func (InputMediaVideo) inputMedia() {}

// SendMediaGroup send a group of photos or videos as an album
func (c *Client) SendMediaGroup(chatID string, media []InputMedia, opts ...sendOption) ([]*Message, error) {
	req := url.Values{}
	req.Set("chat_id", chatID)
	m, _ := json.Marshal(media)
	req.Set("media", string(m))
	for _, opt := range opts {
		opt(req)
	}
	var msgs []*Message
	err := c.doRequest("sendMediaGroup", req, &msgs)
	return msgs, err
}

// SendLocation options
var (
	OptLivePeriod = func(period int) sendOption {
		return func(v url.Values) {
			v.Set("live_period", fmt.Sprint(period))
		}
	}
)

/*
SendLocation sends point on the map to chat. Available options:
	- OptLivePeriod(period int)
	- OptDisableNotification
	- OptReplyToMessageID(id int)
	- OptInlineKeyboardMarkup(markup *InlineKeyboardMarkup)
	- OptReplyKeyboardMarkup(markup *ReplyKeyboardMarkup)
	- OptReplyKeyboardRemove
	- OptReplyKeyboardRemoveSelective
	- OptForceReply
	- OptForceReplySelective
*/
func (c *Client) SendLocation(chatID string, latitude, longitude float64, opts ...sendOption) (*Message, error) {
	req := url.Values{}
	req.Set("chat_id", chatID)
	req.Set("latitude", fmt.Sprint(latitude))
	req.Set("longitude", fmt.Sprint(longitude))
	for _, opt := range opts {
		opt(req)
	}
	msg := &Message{}
	err := c.doRequest("sendLocation", req, msg)
	return msg, err
}

/*
EditMessageLiveLocation edits location in message sent by the bot. Available options:
	- OptInlineKeyboardMarkup(markup *InlineKeyboardMarkup)
*/
func (c *Client) EditMessageLiveLocation(chatID string, messageID int, latitude, longitude float64, opts ...sendOption) (*Message, error) {
	req := url.Values{}
	req.Set("chat_id", chatID)
	req.Set("message_id", fmt.Sprint(messageID))
	req.Set("latitude", fmt.Sprint(latitude))
	req.Set("longitude", fmt.Sprint(longitude))
	for _, opt := range opts {
		opt(req)
	}
	msg := &Message{}
	err := c.doRequest("editMessageLiveLocation", req, msg)
	return msg, err
}

/*
EditInlineMessageLiveLocation edits location in message sent via the bot (using inline mode). Available options:
	- OptInlineKeyboardMarkup(markup *InlineKeyboardMarkup)
*/
func (c *Client) EditInlineMessageLiveLocation(inlineMessageID string, latitude, longitude float64, opts ...sendOption) error {
	req := url.Values{}
	req.Set("inline_message_id", inlineMessageID)
	req.Set("latitude", fmt.Sprint(latitude))
	req.Set("longitude", fmt.Sprint(longitude))
	for _, opt := range opts {
		opt(req)
	}
	var edited bool
	err := c.doRequest("editMessageLiveLocation", req, &edited)
	return err
}

/*
StopMessageLiveLocation stop updating a live location message sent by the bot. Available options:
	- OptInlineKeyboardMarkup(markup *InlineKeyboardMarkup)
*/
func (c *Client) StopMessageLiveLocation(chatID string, messageID int, opts ...sendOption) (*Message, error) {
	req := url.Values{}
	req.Set("chat_id", chatID)
	req.Set("message_id", fmt.Sprint(messageID))
	for _, opt := range opts {
		opt(req)
	}
	msg := &Message{}
	err := c.doRequest("stopMessageLiveLocation", req, msg)
	return msg, err
}

/*
StopInlineMessageLiveLocation stop updating a live location message sent via the bot (using inline mode). Available options:
	- OptInlineKeyboardMarkup(markup *InlineKeyboardMarkup)
*/
func (c *Client) StopInlineMessageLiveLocation(inlineMessageID string, opts ...sendOption) error {
	req := url.Values{}
	req.Set("inline_message_id", inlineMessageID)
	for _, opt := range opts {
		opt(req)
	}
	var stopped bool
	return c.doRequest("stopMessageLiveLocation", req, &stopped)
}

// SendVenue options
var (
	OptFoursquareID = func(foursquareID string) sendOption {
		return func(v url.Values) {
			v.Set("foursquare_id", foursquareID)
		}
	}
	OptFoursquareType = func(foursquareType string) sendOption {
		return func(v url.Values) {
			v.Set("foursquare_type", foursquareType)
		}
	}
)

/*
SendVenue sends information about a venue. Available options:
	- OptFoursquareID(foursquareID string)
	- OptFoursquareType(foursquareType string)
	- OptDisableNotification
	- OptReplyToMessageID(id int)
	- OptInlineKeyboardMarkup(markup *InlineKeyboardMarkup)
	- OptReplyKeyboardMarkup(markup *ReplyKeyboardMarkup)
	- OptReplyKeyboardRemove
	- OptReplyKeyboardRemoveSelective
	- OptForceReply
	- OptForceReplySelective
*/
func (c *Client) SendVenue(chatID string, latitude, longitude float64, title, address string, opts ...sendOption) (*Message, error) {
	req := url.Values{}
	req.Set("chat_id", chatID)
	req.Set("latitude", fmt.Sprint(latitude))
	req.Set("longitude", fmt.Sprint(longitude))
	req.Set("title", title)
	req.Set("address", address)
	for _, opt := range opts {
		opt(req)
	}
	msg := &Message{}
	err := c.doRequest("sendVenue", req, msg)
	return msg, err
}

// SendContact options
var (
	OptLastName = func(lastName string) sendOption {
		return func(v url.Values) {
			v.Set("last_name", lastName)
		}
	}
	OptVCard = func(vCard string) sendOption {
		return func(v url.Values) {
			v.Set("vcard", vCard)
		}
	}
)

/*
SendContact sends phone contact. Available options:
	- OptLastName(lastName string)
	- OptVCard(vCard string) TODO: implement vCard support (https://tools.ietf.org/html/rfc6350)
	- OptDisableNotification
	- OptReplyToMessageID(id int)
	- OptInlineKeyboardMarkup(markup *InlineKeyboardMarkup)
	- OptReplyKeyboardMarkup(markup *ReplyKeyboardMarkup)
	- OptReplyKeyboardRemove
	- OptReplyKeyboardRemoveSelective
	- OptForceReply
	- OptForceReplySelective
*/
func (c *Client) SendContact(chatID, phoneNumber, firstName string, opts ...sendOption) (*Message, error) {
	req := url.Values{}
	req.Set("chat_id", chatID)
	req.Set("phone_number", phoneNumber)
	req.Set("first_name", firstName)
	for _, opt := range opts {
		opt(req)
	}
	msg := &Message{}
	err := c.doRequest("sendContact", req, msg)
	return msg, err
}

type chatAction string

// Actions for SendChatAction
const (
	ActionTyping          chatAction = "typing"
	ActionUploadPhoto     chatAction = "upload_photo"
	ActionRecordVideo     chatAction = "record_video"
	ActionUploadVideo     chatAction = "upload_video"
	ActionRecordAudio     chatAction = "record_audio"
	ActionUploadAudio     chatAction = "upload_audio"
	ActionUploadDocument  chatAction = "upload_document"
	ActionFindLocation    chatAction = "find_location"
	ActionRecordVideoNote chatAction = "record_video_note"
	ActionUploadVideoNote chatAction = "upload_video_note"
)

/*
SendChatAction sends bot chat action. Available actions:
	- ActionTyping
	- ActionUploadPhoto
	- ActionRecordVideo
	- ActionUploadVideo
	- ActionRecordAudio
	- ActionUploadAudio
	- ActionUploadDocument
	- ActionFindLocation
	- ActionRecordVideoNote
	- ActionUploadVideoNote
*/
func (c *Client) SendChatAction(chatID string, action chatAction) error {
	req := url.Values{}
	req.Set("chat_id", chatID)
	req.Set("action", string(action))
	var sent bool
	return c.doRequest("sendChatAction", req, &sent)
}

// UserProfilePhotos represent a user's profile pictures
type UserProfilePhotos struct {
	TotalCount int           `json:"total_count"`
	Photos     [][]PhotoSize `json:"photos"`
}

// GetUserProfilePhotos options
var (
	OptOffset = func(offset int) sendOption {
		return func(v url.Values) {
			v.Set("offset", fmt.Sprint(offset))
		}
	}
	OptLimit = func(limit int) sendOption {
		return func(v url.Values) {
			v.Set("limit", fmt.Sprint(limit))
		}
	}
)

/*
GetUserProfilePhotos returs user's profile pictures. Available options:
	- OptOffset(offset int)
	- OptLimit(limit int)
*/
func (c *Client) GetUserProfilePhotos(userID int, opts ...sendOption) (*UserProfilePhotos, error) {
	req := url.Values{}
	req.Set("user_id", fmt.Sprint(userID))
	for _, opt := range opts {
		opt(req)
	}
	photos := &UserProfilePhotos{}
	err := c.doRequest("getUserProfilePhotos", req, photos)
	return photos, err
}

// File object represents a file ready to be downloaded
type File struct {
	FileID   string `json:"file_id"`
	FileSize int    `json:"file_size"`
	FilePath string `json:"file_path"` // use https://api.telegram.org/file/bot<token>/<file_path> to download
}

/*
GetFile returns File object by fileID.
*/
func (c *Client) GetFile(fileID string) (*File, error) {
	req := url.Values{}
	req.Set("file_id", fileID)
	file := &File{}
	err := c.doRequest("getFile", req, file)
	return file, err
}

// KickChatMember options
var (
	OptUntilDate = func(date time.Time) sendOption {
		return func(v url.Values) {
			v.Set("until_date", fmt.Sprint(date.Unix()))
		}
	}
)

/*
KickChatMember kicks user from group, supergroup or channel. Available options:
	- OptUntilDate(date time.Time)
*/
func (c *Client) KickChatMember(chatID string, userID int, opts ...sendOption) error {
	req := url.Values{}
	req.Set("chat_id", chatID)
	req.Set("user_id", fmt.Sprint(userID))
	for _, opt := range opts {
		opt(req)
	}
	var kicked bool
	return c.doRequest("kickChatMember", req, &kicked)
}

/*
UnbanChatMember unban a previously kicked user in a supergroup or channel
*/
func (c *Client) UnbanChatMember(chatID string, userID int) error {
	req := url.Values{}
	req.Set("chat_id", chatID)
	req.Set("user_id", fmt.Sprint(userID))
	var unbanned bool
	return c.doRequest("unbanChatMember", req, &unbanned)
}

// Restrictions for user in supergroup
type Restrictions struct {
	CanSendMessages       bool
	CanSendMediaMessages  bool
	CanSendOtherMessages  bool
	CanAddWebPagePreviews bool
}

/*
RestrictChatMember restrict a user in a supergroup. Available options:
	- OptUntilDate(date time.Time)
*/
func (c *Client) RestrictChatMember(chatID string, userID int, r *Restrictions, opts ...sendOption) error {
	req := url.Values{}
	req.Set("chat_id", chatID)
	req.Set("user_id", fmt.Sprint(userID))
	req.Set("can_send_messages", fmt.Sprint(r.CanSendMessages))
	req.Set("can_send_media_messages", fmt.Sprint(r.CanSendMediaMessages))
	req.Set("can_send_other_messages", fmt.Sprint(r.CanSendOtherMessages))
	req.Set("can_add_web_page_previews", fmt.Sprint(r.CanAddWebPagePreviews))
	for _, opt := range opts {
		opt(req)
	}
	var restricted bool
	return c.doRequest("restrictChatMember", req, &restricted)
}

// Promotions give user permitions in a supergroup or channel.
type Promotions struct {
	CanChangeInfo      bool
	CanPostMessages    bool
	CanEditMessages    bool
	CanDeleteMessages  bool
	CanInviteUsers     bool
	CanRestrictMembers bool
	CanPinMessages     bool
	CanPromoteMembers  bool
}

/*
PromoteChatMember promote or demote a user in a supergroup or a channel
*/
func (c *Client) PromoteChatMember(chatID string, userID int, p *Promotions) error {
	req := url.Values{}
	req.Set("chat_id", chatID)
	req.Set("user_id", fmt.Sprint(userID))
	req.Set("can_change_info", fmt.Sprint(p.CanChangeInfo))
	req.Set("can_post_messages", fmt.Sprint(p.CanPostMessages))
	req.Set("can_edit_messages", fmt.Sprint(p.CanEditMessages))
	req.Set("can_delete_messages", fmt.Sprint(p.CanDeleteMessages))
	req.Set("can_invite_users", fmt.Sprint(p.CanInviteUsers))
	req.Set("can_restrict_members", fmt.Sprint(p.CanRestrictMembers))
	req.Set("can_pin_messages", fmt.Sprint(p.CanPinMessages))
	req.Set("can_promote_member", fmt.Sprint(p.CanPromoteMembers))
	var promoted bool
	return c.doRequest("promoteChatMember", req, &promoted)
}

/*
ExportChatInviteLink generate a new invite link for a chat; any previously generated link is revoked
*/
func (c *Client) ExportChatInviteLink(chatID string) (string, error) {
	req := url.Values{}
	req.Set("chat_id", chatID)
	var link string
	err := c.doRequest("exportChatInviteLink", req, &link)
	return link, err
}

/*
SetChatPhoto set a new profile photo for the chat
*/
func (c *Client) SetChatPhoto(chatID string, filename string) error {
	req := url.Values{}
	req.Set("chat_id", chatID)
	var updated bool
	return c.doRequestWithFiles("setChatPhoto", req, &updated, inputFile{field: "photo", name: filename})
}

/*
DeleteChatPhoto deleta a chat photo
*/
func (c *Client) DeleteChatPhoto(chatID string) error {
	req := url.Values{}
	req.Set("chat_id", chatID)
	var deleted bool
	return c.doRequest("deleteChatPhoto", req, &deleted)
}

/*
SetChatTitle change the title of the chat
*/
func (c *Client) SetChatTitle(chatID, title string) error {
	req := url.Values{}
	req.Set("chat_id", chatID)
	req.Set("title", title)
	var set bool
	return c.doRequest("setChatTitle", req, &set)
}

/*
SetChatDescription change the description of a supergroup or a channel
*/
func (c *Client) SetChatDescription(chatID, description string) error {
	req := url.Values{}
	req.Set("chat_id", chatID)
	req.Set("description", description)
	var set bool
	return c.doRequest("setChatDescription", req, &set)
}

/*
PinChatMessage pin a message in a supergroup or a channel. Available options:
	- OptDisableNotification
*/
func (c *Client) PinChatMessage(chatID string, messageID int, opts ...sendOption) error {
	req := url.Values{}
	req.Set("chat_id", chatID)
	req.Set("message_id", fmt.Sprint(messageID))
	for _, opt := range opts {
		opt(req)
	}
	var pinned bool
	return c.doRequest("pinChatMessage", req, &pinned)
}

/*
UnpinChatMessage unpin a message in a supergroup or a channel
*/
func (c *Client) UnpinChatMessage(chatID string) error {
	req := url.Values{}
	req.Set("chat_id", chatID)
	var unpinned bool
	return c.doRequest("unpinChatMessage", req, &unpinned)
}

/*
LeaveChat leave a group, supergroup or channel
*/
func (c *Client) LeaveChat(chatID string) error {
	req := url.Values{}
	req.Set("chat_id", chatID)
	var left bool
	return c.doRequest("leaveChat", req, &left)
}

/*
GetChat get up to date information about the chat
*/
func (c *Client) GetChat(chatID string) (*Chat, error) {
	req := url.Values{}
	req.Set("chat_id", chatID)
	chat := &Chat{}
	err := c.doRequest("getChat", req, chat)
	return chat, err
}

// ChatMember contains information about one member of a chat
type ChatMember struct {
	User                  User   `json:"user"`
	Status                string `json:"status"`
	UntilDate             int    `json:"until_date"`
	CanBeEdited           bool   `json:"can_be_edited"`
	CanChangeInfo         bool   `json:"can_change_info"`
	CanPostMessages       bool   `json:"can_post_messages"`
	CanEditMessages       bool   `json:"can_edit_messages"`
	CanDeleteMessages     bool   `json:"can_delete_messages"`
	CanInviteUsers        bool   `json:"can_invite_users"`
	CanRestrictMembers    bool   `json:"can_restrict_members"`
	CanPinMessages        bool   `json:"can_pin_messages"`
	CanPromoteMembers     bool   `json:"can_promote_members"`
	IsMember              bool   `json:"is_member"`
	CanSendMessages       bool   `json:"can_send_messages"`
	CanSendMediaMessages  bool   `json:"can_send_media_messages"`
	CanSendOtherMessages  bool   `json:"can_send_other_messages"`
	CanAddWebPagePreviews bool   `json:"can_add_web_page_previews"`
}

/*
GetChatAdministrators get a list of administrators in a chat
*/
func (c *Client) GetChatAdministrators(chatID string) ([]*ChatMember, error) {
	req := url.Values{}
	req.Set("chat_id", chatID)
	members := []*ChatMember{}
	err := c.doRequest("getChatAdministrators", req, &members)
	return members, err
}

/*
GetChatMembersCount returns the number of members in chat
*/
func (c *Client) GetChatMembersCount(chatID string) (int, error) {
	req := url.Values{}
	req.Set("chat_id", chatID)
	var count int
	err := c.doRequest("getChatMembersCount", req, &count)
	return count, err
}

/*
GetChatMember get information about a member of a chat
*/
func (c *Client) GetChatMember(chatID string, userID int) (*ChatMember, error) {
	req := url.Values{}
	req.Set("chat_id", chatID)
	req.Set("user_id", fmt.Sprint(userID))
	member := &ChatMember{}
	err := c.doRequest("getChatMember", req, member)
	return member, err
}

/*
SetChatStickerSet set a new group sticker set for a supergroup
*/
func (c *Client) SetChatStickerSet(chatID, stickerSetName string) error {
	req := url.Values{}
	req.Set("chat_id", chatID)
	req.Set("sticker_set_name", stickerSetName)
	var set bool
	return c.doRequest("setChatStickerSet", req, &set)
}

/*
DeleteChatStickerSet delete a group sticker set from a supergroup
*/
func (c *Client) DeleteChatStickerSet(chatID string) error {
	req := url.Values{}
	req.Set("chat_id", chatID)
	var deleted bool
	return c.doRequest("deleteChatStickerSet", req, &deleted)
}

// Options for AnswerCallbackQuery
var (
	OptText = func(text string) sendOption {
		return func(v url.Values) {
			v.Set("text", text)
		}
	}
	OptShowAlert = func(v url.Values) {
		v.Set("show_alert", "true")
	}
	OptURL = func(u string) sendOption {
		return func(v url.Values) {
			v.Set("url", u)
		}
	}
	OptCacheTime = func(d time.Duration) sendOption {
		return func(v url.Values) {
			v.Set("cache_time", fmt.Sprint(int(d.Seconds())))
		}
	}
)

/*
AnswerCallbackQuery send answer to callback query sent from inline keyboard. Available options:
	- OptText(text string)
	- OptShowAlert
	- OptURL(url string)
	- OptCacheTime(d time.Duration)
*/
func (c *Client) AnswerCallbackQuery(callbackQueryID string, opts ...sendOption) error {
	req := url.Values{}
	req.Set("callback_query_id", callbackQueryID)
	for _, opt := range opts {
		opt(req)
	}
	var success bool
	return c.doRequest("answerCallbackQuery", req, &success)
}

/*
EditMessageText edit text and game messages sent by the bot. Available options:
	- OptParseModeHTML
	- OptParseModeMarkdown
	- OptDisableWebPagePreview
	- OptInlineKeyboardMarkup(markup *InlineKeyboardMarkup)
*/
func (c *Client) EditMessageText(chatID string, messageID int, text string, opts ...sendOption) (*Message, error) {
	req := url.Values{}
	req.Set("chat_id", chatID)
	req.Set("message_id", fmt.Sprint(messageID))
	req.Set("text", text)
	for _, opt := range opts {
		opt(req)
	}
	msg := &Message{}
	err := c.doRequest("editMessageText", req, msg)
	return msg, err
}

/*
EditInlineMessageText edit text and game messages sent via the bot (for inline bots). Available options:
	- OptParseModeHTML
	- OptParseModeMarkdown
	- OptDisableWebPagePreview
	- OptInlineKeyboardMarkup(markup *InlineKeyboardMarkup)
*/
func (c *Client) EditInlineMessageText(inlineMessageID, text string, opts ...sendOption) error {
	req := url.Values{}
	req.Set("inline_message_id", inlineMessageID)
	req.Set("text", text)
	for _, opt := range opts {
		opt(req)
	}
	var edited bool
	return c.doRequest("editMessageText", req, &edited)
}

/*
EditMessageCaption edit message caption sent by the bot. Available options:
	- OptParseModeHTML
	- OptParseModeMarkdown
	- OptInlineKeyboardMarkup(markup *InlineKeyboardMarkup)
*/
func (c *Client) EditMessageCaption(chatID string, messageID int, caption string, opts ...sendOption) (*Message, error) {
	req := url.Values{}
	req.Set("chat_id", chatID)
	req.Set("message_id", fmt.Sprint(messageID))
	req.Set("caption", caption)
	for _, opt := range opts {
		opt(req)
	}
	msg := &Message{}
	err := c.doRequest("editMessageCaption", req, msg)
	return msg, err
}

/*
EditInlineMessageCaption edit message caption sent via the bot (for inline bots). Available options:
	- OptParseModeHTML
	- OptParseModeMarkdown
	- OptInlineKeyboardMarkup(markup *InlineKeyboardMarkup)
*/
func (c *Client) EditInlineMessageCaption(inlineMessageID, caption string, opts ...sendOption) error {
	req := url.Values{}
	req.Set("inline_message_id", inlineMessageID)
	req.Set("caption", caption)
	for _, opt := range opts {
		opt(req)
	}
	var edited bool
	return c.doRequest("editMessageCaption", req, &edited)
}

/*
EditMessageReplyMarkup edit only the reply markup of messages sent by the bot. Available options:
	- OptInlineKeyboardMarkup(markup *InlineKeyboardMarkup)
*/
func (c *Client) EditMessageReplyMarkup(chatID string, messageID int, opts ...sendOption) (*Message, error) {
	req := url.Values{}
	req.Set("chat_id", chatID)
	req.Set("message_id", fmt.Sprint(messageID))
	for _, opt := range opts {
		opt(req)
	}
	msg := &Message{}
	err := c.doRequest("editMessageReplyMarkup", req, msg)
	return msg, err
}

/*
EditInlineMessageReplyMarkup edit only the reply markup of messages sent by the bot. Available options:
	- OptInlineKeyboardMarkup(markup *InlineKeyboardMarkup)
*/
func (c *Client) EditInlineMessageReplyMarkup(inlineMessageID string, opts ...sendOption) error {
	req := url.Values{}
	req.Set("inline_message_id", inlineMessageID)
	for _, opt := range opts {
		opt(req)
	}
	var edited bool
	return c.doRequest("editMessageReplyMarkup", req, &edited)
}

/*
DeleteMessage delete a message, including service messages
*/
func (c *Client) DeleteMessage(chatID string, messageID int) error {
	req := url.Values{}
	req.Set("chat_id", chatID)
	req.Set("message_id", fmt.Sprint(messageID))
	var deleted bool
	return c.doRequest("deleteMessage", req, &deleted)
}

/*
SendStickerFile send .webp file sticker. Available options:
	- OptDisableNotification
	- OptReplyToMessageID(id int)
	- OptInlineKeyboardMarkup(markup *InlineKeyboardMarkup)
	- OptReplyKeyboardMarkup(markup *ReplyKeyboardMarkup)
	- OptReplyKeyboardRemove
	- OptReplyKeyboardRemoveSelective
	- OptForceReply
	- OptForceReplySelective
*/
func (c *Client) SendStickerFile(chatID string, filename string, opts ...sendOption) (*Message, error) {
	req := url.Values{}
	req.Set("chat_id", chatID)
	for _, opt := range opts {
		opt(req)
	}
	msg := &Message{}
	err := c.doRequestWithFiles("sendSticker", req, msg, inputFile{field: "sticker", name: filename})
	return msg, err
}

/*
SendSticker send previously uploaded sticker. Available options:
	- OptDisableNotification
	- OptReplyToMessageID(id int)
	- OptInlineKeyboardMarkup(markup *InlineKeyboardMarkup)
	- OptReplyKeyboardMarkup(markup *ReplyKeyboardMarkup)
	- OptReplyKeyboardRemove
	- OptReplyKeyboardRemoveSelective
	- OptForceReply
	- OptForceReplySelective
*/
func (c *Client) SendSticker(chatID, fileID string, opts ...sendOption) (*Message, error) {
	req := url.Values{}
	req.Set("chat_id", chatID)
	req.Set("sticker", fileID)
	for _, opt := range opts {
		opt(req)
	}
	msg := &Message{}
	err := c.doRequest("sendSticker", req, msg)
	return msg, err
}

// StickerSet represents sticker set
type StickerSet struct {
	Name          string    `json:"name"`
	Title         string    `json:"title"`
	IsAnimated    bool      `json:"is_animated"`
	ContainsMasks bool      `json:"contains_masks"`
	Stickers      []Sticker `json:"stickers"`
}

/*
GetStickerSet get a sticker set
*/
func (c *Client) GetStickerSet(name string) (*StickerSet, error) {
	req := url.Values{}
	req.Set("name", name)
	set := &StickerSet{}
	err := c.doRequest("getStickerSet", req, set)
	return set, err
}

/*
UploadStickerFile upload a .png file with a sticker for later use in CreateNewStickerSet and AddStickerToSet
*/
func (c *Client) UploadStickerFile(userID int, filename string) (*File, error) {
	req := url.Values{}
	req.Set("user_id", fmt.Sprint(userID))
	file := &File{}
	err := c.doRequestWithFiles("uploadStickerFile", req, &file, inputFile{field: "png_sticker", name: filename})
	return file, err
}

// CreateNewStickerSet options
var (
	OptContainsMasks = func(v url.Values) {
		v.Set("contains_masks", "true")
	}
	OptMaskPosition = func(pos *MaskPosition) sendOption {
		return func(v url.Values) {
			p, _ := json.Marshal(pos)
			v.Set("mask_position", string(p))
		}
	}
)

/*
CreateNewStickerSetFile creates new sticker set with sticker file. Available options:
	- OptContainsMasks
	- OptMaskPosition(pos *MaskPosition)
*/
func (c *Client) CreateNewStickerSetFile(userID int, name, title, stickerFilename, emojis string, opts ...sendOption) error {
	req := url.Values{}
	req.Set("user_id", fmt.Sprint(userID))
	req.Set("name", name)
	req.Set("title", title)
	req.Set("emojis", emojis)
	for _, opt := range opts {
		opt(req)
	}
	var created bool
	return c.doRequestWithFiles("createNewStickerSet", req, &created, inputFile{field: "png_sticker", name: stickerFilename})
}

/*
CreateNewStickerSet creates new sticker set with previously uploaded file. Available options:
	- OptContainsMasks
	- OptMaskPosition(pos *MaskPosition)
*/
func (c *Client) CreateNewStickerSet(userID int, name, title, fileID, emojis string, opts ...sendOption) error {
	req := url.Values{}
	req.Set("user_id", fmt.Sprint(userID))
	req.Set("name", name)
	req.Set("title", title)
	req.Set("png_sticker", fileID)
	req.Set("emojis", emojis)
	for _, opt := range opts {
		opt(req)
	}
	var created bool
	return c.doRequest("createNewStickerSet", req, &created)
}

/*
AddStickerToSetFile add a new sticker file to a set created by the bot. Available options:
	- OptMaskPosition(pos *MaskPosition)
*/
func (c *Client) AddStickerToSetFile(userID int, name, filename, emojis string, opts ...sendOption) error {
	req := url.Values{}
	req.Set("user_id", fmt.Sprint(userID))
	req.Set("name", name)
	req.Set("emojis", emojis)
	for _, opt := range opts {
		opt(req)
	}
	var added bool
	return c.doRequestWithFiles("addStickerToSet", req, &added, inputFile{field: "png_sticker", name: filename})
}

/*
AddStickerToSet add a new sticker to a set created by the bot. Available options:
	- OptMaskPosition(pos *MaskPosition)
*/
func (c *Client) AddStickerToSet(userID int, name, fileID, emojis string, opts ...sendOption) error {
	req := url.Values{}
	req.Set("user_id", fmt.Sprint(userID))
	req.Set("name", name)
	req.Set("png_sticker", fileID)
	req.Set("emojis", emojis)
	for _, opt := range opts {
		opt(req)
	}
	var added bool
	return c.doRequestWithFiles("addStickerToSet", req, &added)
}

/*
SetStickerPositionInSet move a sticker in a set created by the bot to a specific position
*/
func (c *Client) SetStickerPositionInSet(fileID string, pos int) error {
	req := url.Values{}
	req.Set("sticker", fileID)
	req.Set("position", fmt.Sprint(pos))
	var set bool
	return c.doRequest("setStickerPositionInSet", req, &set)
}

/*
DeleteStickerFromSet delete a sticker from a set created by the bot
*/
func (c *Client) DeleteStickerFromSet(fileID string) error {
	req := url.Values{}
	req.Set("sticker", fileID)
	var deleted bool
	return c.doRequest("deleteStickerFromSet", req, &deleted)
}

// InputMessageContent content of a message to be sent as a result of an inline query
type InputMessageContent interface {
	inputMessageContent()
}

var (
	_ InputMessageContent = InputTextMessageContent{}
	_ InputMessageContent = InputLocationMessageContent{}
	_ InputMessageContent = InputVenueMessageContent{}
	_ InputMessageContent = InputContactMessageContent{}
)

// InputTextMessageContent represents the content of a text message to be sent as the result of an inline query
type InputTextMessageContent struct {
	MessageText           string `json:"message_text"`
	ParseMode             string `json:"parse_mode"`
	DisableWebPagePreview bool   `json:"disable_web_page_preview"`
}

func (InputTextMessageContent) inputMessageContent() {}

// InputLocationMessageContent represents the content of a location message to be sent as the result of an inline query
type InputLocationMessageContent struct {
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	LivePeriod int     `json:"live_period"`
}

func (InputLocationMessageContent) inputMessageContent() {}

// InputVenueMessageContent represents the content of a venue message to be sent as the result of an inline query
type InputVenueMessageContent struct {
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	Title          string  `json:"title"`
	Address        string  `json:"address"`
	FoursquareID   string  `json:"foursquare_id"`
	FoursquareType string  `json:"foursquare_type"`
}

func (InputVenueMessageContent) inputMessageContent() {}

// InputContactMessageContent represents the content of a contact message to be sent as the result of an inline query
type InputContactMessageContent struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	VCard       string `json:"vcard"`
}

func (InputContactMessageContent) inputMessageContent() {}

// InlineQueryResult represents one result of an inline query
type InlineQueryResult interface {
	inlineQueryResult()
}

var (
	_ InlineQueryResult = InlineQueryResultArticle{}
	_ InlineQueryResult = InlineQueryResultPhoto{}
	_ InlineQueryResult = InlineQueryResultGif{}
	_ InlineQueryResult = InlineQueryResultMpeg4Gif{}
	_ InlineQueryResult = InlineQueryResultVideo{}
	_ InlineQueryResult = InlineQueryResultAudio{}
	_ InlineQueryResult = InlineQueryResultVoice{}
	_ InlineQueryResult = InlineQueryResultDocument{}
	_ InlineQueryResult = InlineQueryResultLocation{}
	_ InlineQueryResult = InlineQueryResultVenue{}
	_ InlineQueryResult = InlineQueryResultContact{}
	_ InlineQueryResult = InlineQueryResultGame{}
	_ InlineQueryResult = InlineQueryResultCachedPhoto{}
	_ InlineQueryResult = InlineQueryResultCachedGif{}
	_ InlineQueryResult = InlineQueryResultCachedMpeg4Gif{}
	_ InlineQueryResult = InlineQueryResultCachedSticker{}
	_ InlineQueryResult = InlineQueryResultCachedDocument{}
	_ InlineQueryResult = InlineQueryResultCachedVideo{}
	_ InlineQueryResult = InlineQueryResultCachedVoice{}
	_ InlineQueryResult = InlineQueryResultCachedAudio{}
)

// InlineQueryResultArticle link to an article or web page
type InlineQueryResultArticle struct {
	Type                string                `json:"type"`
	ID                  string                `json:"id"`
	Title               string                `json:"title"`
	InputMessageContent InputMessageContent   `json:"input_message_content"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	URL                 string                `json:"url,omitempty"`
	HideURL             bool                  `json:"hide_url,omitempty"`
	Description         string                `json:"description,omitempty"`
	ThumbURL            string                `json:"thumb_url,omitempty"`
	ThumbWidth          int                   `json:"thumb_width,omitempty"`
	ThumbHeight         int                   `json:"thumb_height,omitempty"`
}

func (InlineQueryResultArticle) inlineQueryResult() {}

// InlineQueryResultPhoto represents a link to a photo
type InlineQueryResultPhoto struct {
	Type                string                `json:"type"`
	ID                  string                `json:"id"`
	PhotoURL            string                `json:"photo_url"`
	ThumbURL            string                `json:"thumb_url"`
	PhotoWidth          int                   `json:"photo_width,omitempty"`
	PhotoHeight         int                   `json:"photo_height,omitempty"`
	Title               string                `json:"title,omitempty"`
	Description         string                `json:"description,omitempty"`
	Caption             string                `json:"caption,omitempty"`
	ParseMode           string                `json:"parse_mode,omitempty"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent *InputMessageContent  `json:"input_message_content,omitempty"`
}

func (InlineQueryResultPhoto) inlineQueryResult() {}

// InlineQueryResultGif represents a link to an animated GIF file
type InlineQueryResultGif struct {
	Type                string                `json:"type"`
	ID                  string                `json:"id"`
	GifURL              string                `json:"gif_url"`
	GifWidth            int                   `json:"gif_width,omitempty"`
	GifHeight           int                   `json:"gif_height,omitempty"`
	GifDuration         int                   `json:"gif_duration,omitempty"`
	ThumbURL            string                `json:"thumb_url"`
	Title               string                `json:"title,omitempty"`
	Caption             string                `json:"caption,omitempty"`
	ParseMode           string                `json:"parse_mode,omitempty"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent *InputMessageContent  `json:"input_message_content,omitempty"`
}

func (InlineQueryResultGif) inlineQueryResult() {}

// InlineQueryResultMpeg4Gif represents a link to a video animation (H.264/MPEG-4 AVC video without sound)
type InlineQueryResultMpeg4Gif struct {
	Type                string                `json:"type"`
	ID                  string                `json:"id"`
	Mpeg4URL            string                `json:"mpeg4_url"`
	Mpeg4Width          int                   `json:"mpeg4_width,omitempty"`
	Mpeg4Height         int                   `json:"mpeg4_height,omitempty"`
	Mpeg4Duration       int                   `json:"mpeg4_duration,omitempty"`
	ThumbURL            string                `json:"thumb_url"`
	Title               string                `json:"title,omitempty"`
	Caption             string                `json:"caption,omitempty"`
	ParseMode           string                `json:"parse_mode,omitempty"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent *InputMessageContent  `json:"input_message_content,omitempty"`
}

func (InlineQueryResultMpeg4Gif) inlineQueryResult() {}

// InlineQueryResultVideo represents a link to a page containing an embedded video player or a video file
type InlineQueryResultVideo struct {
	Type                string                `json:"type"`
	ID                  string                `json:"id"`
	VideoURL            string                `json:"video_url"`
	MimeType            string                `json:"mime_type"`
	ThumbURL            string                `json:"thumb_url"`
	Title               string                `json:"title"`
	Caption             string                `json:"caption,omitempty"`
	ParseMode           string                `json:"parse_mode,omitempty"`
	VideoWidth          int                   `json:"video_width,omitempty"`
	Videoheight         int                   `json:"video_height,omitempty"`
	VideoDuration       int                   `json:"video_duration,omitempty"`
	Description         string                `json:"description,omitempty"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent *InputMessageContent  `json:"input_message_content,omitempty"`
}

func (InlineQueryResultVideo) inlineQueryResult() {}

// InlineQueryResultAudio represents a link to an mp3 audio file
type InlineQueryResultAudio struct {
	Type                string                `json:"type"`
	ID                  string                `json:"id"`
	AudioURL            string                `json:"audio_url"`
	Title               string                `json:"title"`
	Caption             string                `json:"caption,omitempty"`
	ParseMode           string                `json:"parse_mode,omitempty"`
	Performer           string                `json:"performer,omitempty"`
	AudioDuration       int                   `json:"audio_duration,omitempty"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent *InputMessageContent  `json:"input_message_content,omitempty"`
}

func (InlineQueryResultAudio) inlineQueryResult() {}

// InlineQueryResultVoice represents a link to a voice recording in an .ogg container encoded with OPUS
type InlineQueryResultVoice struct {
	Type                string                `json:"type"`
	ID                  string                `json:"id"`
	VoiceURL            string                `json:"voice_url"`
	Title               string                `json:"title"`
	Caption             string                `json:"caption,omitempty"`
	ParseMode           string                `json:"parse_mode,omitempty"`
	Performer           string                `json:"performer,omitempty"`
	VoiceDuration       int                   `json:"voice_duration,omitempty"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent *InputMessageContent  `json:"input_message_content,omitempty"`
}

func (InlineQueryResultVoice) inlineQueryResult() {}

// InlineQueryResultDocument represents a link to a file
type InlineQueryResultDocument struct {
	Type                string                `json:"type"`
	ID                  string                `json:"id"`
	Title               string                `json:"title"`
	Caption             string                `json:"caption"`
	ParseMode           string                `json:"parse_mode"`
	DocumentURL         string                `json:"document_url"`
	MimeType            string                `json:"mime_type"`
	Description         string                `json:"description,omitempty"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent *InputMessageContent  `json:"input_message_content,omitempty"`
	ThumbURL            string                `json:"thumb_url,omitempty"`
	ThumbWidth          int                   `json:"thumb_width,omitempty"`
	ThumbHeight         int                   `json:"thumb_height,omitempty"`
}

func (InlineQueryResultDocument) inlineQueryResult() {}

// InlineQueryResultLocation represents a location on a map
type InlineQueryResultLocation struct {
	Type                string                `json:"type"`
	ID                  string                `json:"id"`
	Latitude            float64               `json:"latitude"`
	Longitude           float64               `json:"longitude"`
	Title               string                `json:"title"`
	LivePeriod          int                   `json:"live_period,omitempty"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent *InputMessageContent  `json:"input_message_content,omitempty"`
	ThumbURL            string                `json:"thumb_url,omitempty"`
	ThumbWidth          int                   `json:"thumb_width,omitempty"`
	ThumbHeight         int                   `json:"thumb_height,omitempty"`
}

func (InlineQueryResultLocation) inlineQueryResult() {}

// InlineQueryResultVenue represents a venue
type InlineQueryResultVenue struct {
	Type                string                `json:"type"`
	ID                  string                `json:"id"`
	Latitude            float64               `json:"latitude"`
	Longitude           float64               `json:"longitude"`
	Title               string                `json:"title"`
	Address             string                `json:"address"`
	FoursquareID        string                `json:"foursquare_id,omitempty"`
	FoursquareType      string                `json:"foursquare_type,omitempty"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent *InputMessageContent  `json:"input_message_content,omitempty"`
	ThumbURL            string                `json:"thumb_url,omitempty"`
	ThumbWidth          int                   `json:"thumb_width,omitempty"`
	ThumbHeight         int                   `json:"thumb_height,omitempty"`
}

func (InlineQueryResultVenue) inlineQueryResult() {}

// InlineQueryResultContact represents a contact with a phone number
type InlineQueryResultContact struct {
	Type                string                `json:"type"`
	ID                  string                `json:"id"`
	PhoneNumber         string                `json:"phone_number"`
	FirstName           string                `json:"first_name"`
	LastName            string                `json:"last_name,omitempty"`
	VCard               string                `json:"vcard,omitempty"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent *InputMessageContent  `json:"input_message_content,omitempty"`
	ThumbURL            string                `json:"thumb_url,omitempty"`
	ThumbWidth          int                   `json:"thumb_width,omitempty"`
	ThumbHeight         int                   `json:"thumb_height,omitempty"`
}

func (InlineQueryResultContact) inlineQueryResult() {}

// InlineQueryResultGame represents a Game
type InlineQueryResultGame struct {
	Type          string                `json:"type"`
	ID            string                `json:"id"`
	GameShortName string                `json:"game_short_name"`
	ReplyMarkup   *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

func (InlineQueryResultGame) inlineQueryResult() {}

// InlineQueryResultCachedPhoto represents a link to a photo stored on the Telegram servers
type InlineQueryResultCachedPhoto struct {
	Type                string                `json:"type"`
	ID                  string                `json:"id"`
	PhotoFileID         string                `json:"photo_file_id"`
	Title               string                `json:"title,omitempty"`
	Description         string                `json:"description,omitempty"`
	Caption             string                `json:"caption,omitempty"`
	ParseMode           string                `json:"parse_mode,omitempty"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent *InputMessageContent  `json:"input_message_content,omitempty"`
}

func (InlineQueryResultCachedPhoto) inlineQueryResult() {}

// InlineQueryResultCachedGif represents a link to an animated GIF file stored on the Telegram servers
type InlineQueryResultCachedGif struct {
	Type                string                `json:"type"`
	ID                  string                `json:"id"`
	GifFileID           string                `json:"gif_file_id"`
	Title               string                `json:"title,omitempty"`
	Caption             string                `json:"caption,omitempty"`
	ParseMode           string                `json:"parse_mode,omitempty"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent *InputMessageContent  `json:"input_message_content,omitempty"`
}

func (InlineQueryResultCachedGif) inlineQueryResult() {}

// InlineQueryResultCachedMpeg4Gif represents a link to a video animation (H.264/MPEG-4 AVC video without sound)
// stored on the Telegram servers
type InlineQueryResultCachedMpeg4Gif struct {
	Type                string                `json:"type"`
	ID                  string                `json:"id"`
	Mpeg4FileID         string                `json:"mpeg4_file_id"`
	Title               string                `json:"title,omitempty"`
	Caption             string                `json:"caption,omitempty"`
	ParseMode           string                `json:"parse_mode,omitempty"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent *InputMessageContent  `json:"input_message_content,omitempty"`
}

func (InlineQueryResultCachedMpeg4Gif) inlineQueryResult() {}

// InlineQueryResultCachedSticker represents a link to a sticker stored on the Telegram servers
type InlineQueryResultCachedSticker struct {
	Type                string                `json:"type"`
	ID                  string                `json:"id"`
	StickerFileID       string                `json:"sticker_file_id"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent *InputMessageContent  `json:"input_message_content,omitempty"`
}

func (InlineQueryResultCachedSticker) inlineQueryResult() {}

// InlineQueryResultCachedDocument represents a link to a file
type InlineQueryResultCachedDocument struct {
	Type                string                `json:"type"`
	ID                  string                `json:"id"`
	Title               string                `json:"title"`
	DocumentFileID      string                `json:"document_file_id"`
	Description         string                `json:"description,omitempty"`
	Caption             string                `json:"caption,omitempty"`
	ParseMode           string                `json:"parse_mode,omitempty"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent *InputMessageContent  `json:"input_message_content,omitempty"`
}

func (InlineQueryResultCachedDocument) inlineQueryResult() {}

// InlineQueryResultCachedVideo represents a link to a video file stored on the Telegram servers
type InlineQueryResultCachedVideo struct {
	Type                string                `json:"type"`
	ID                  string                `json:"id"`
	VideoFileID         string                `json:"video_file_id"`
	Title               string                `json:"title"`
	Description         string                `json:"description,omitempty"`
	Caption             string                `json:"caption,omitempty"`
	ParseMode           string                `json:"parse_mode,omitempty"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent *InputMessageContent  `json:"input_message_content,omitempty"`
}

func (InlineQueryResultCachedVideo) inlineQueryResult() {}

// InlineQueryResultCachedVoice represents a link to a voice recording in an .ogg container encoded with OPUS
type InlineQueryResultCachedVoice struct {
	Type                string                `json:"type"`
	ID                  string                `json:"id"`
	VoiceFileID         string                `json:"voice_file_id"`
	Title               string                `json:"title"`
	Caption             string                `json:"caption,omitempty"`
	ParseMode           string                `json:"parse_mode,omitempty"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent *InputMessageContent  `json:"input_message_content,omitempty"`
}

func (InlineQueryResultCachedVoice) inlineQueryResult() {}

// InlineQueryResultCachedAudio represents a link to an mp3 audio file
type InlineQueryResultCachedAudio struct {
	Type                string                `json:"type"`
	ID                  string                `json:"id"`
	AudioFileID         string                `json:"audio_file_id"`
	Caption             string                `json:"caption,omitempty"`
	ParseMode           string                `json:"parse_mode,omitempty"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent *InputMessageContent  `json:"input_message_content,omitempty"`
}

func (InlineQueryResultCachedAudio) inlineQueryResult() {}

// AnswerInlineQuery options
var (
	OptIsPersonal = func(v url.Values) {
		v.Set("is_personal", "true")
	}
	OptNextOffset = func(offset string) sendOption {
		return func(v url.Values) {
			v.Set("next_offset", offset)
		}
	}
	OptSwitchPmText = func(text string) sendOption {
		return func(v url.Values) {
			v.Set("switch_pm_text", text)
		}
	}
	OptSwitchPmParameter = func(param string) sendOption {
		return func(v url.Values) {
			v.Set("switch_pm_parameter", param)
		}
	}
)

/*
AnswerInlineQuery send answer to an inline query. No more than 50 results per query are allowed. Available Options:
	- OptCacheTime(d *time.Duration)
	- OptIsPersonal
	- OptNextOffset(offset string)
	- OptSwitchPmText(text string)
	- OptSwitchPmParameter(param string)
*/
func (c *Client) AnswerInlineQuery(inlineQueryID string, results []InlineQueryResult, opts ...sendOption) error {
	req := url.Values{}
	req.Set("inline_query_id", inlineQueryID)
	res, _ := json.Marshal(results)
	req.Set("results", string(res))
	for _, opt := range opts {
		opt(req)
	}
	var answered bool
	return c.doRequest("answerInlineQuery", req, &answered)
}

// LabeledPrice represents a portion of the price for goods or services
type LabeledPrice struct {
	Label  string `json:"label"`
	Amount int    `json:"amount"`
}

// SendInvoice options
var (
	OptProviderData = func(data string) sendOption {
		return func(v url.Values) {
			v.Set("provider_data", data)
		}
	}
	OptPhotoURL = func(u string) sendOption {
		return func(v url.Values) {
			v.Set("photo_url", u)
		}
	}
	OptPhotoSize = func(size int) sendOption {
		return func(v url.Values) {
			v.Set("photo_size", fmt.Sprint(size))
		}
	}
	OptPhotoWidth = func(width int) sendOption {
		return func(v url.Values) {
			v.Set("photo_width", fmt.Sprint(width))
		}
	}
	OptPhotoHeight = func(height int) sendOption {
		return func(v url.Values) {
			v.Set("photo_height", fmt.Sprint(height))
		}
	}
	OptNeedName                  = func(v url.Values) { v.Set("need_name", "true") }
	OptNeedPhoneNumber           = func(v url.Values) { v.Set("need_phone_number", "true") }
	OptNeedEmail                 = func(v url.Values) { v.Set("need_email", "true") }
	OptNeedShippingAddress       = func(v url.Values) { v.Set("need_shipping_address", "true") }
	OptSendPhoneNumberToProvider = func(v url.Values) { v.Set("send_phone_number_to_provider", "true") }
	OptSendEmailToProvider       = func(v url.Values) { v.Set("send_email_to_provider", "true") }
	OptIsFlexible                = func(v url.Values) { v.Set("is_flexible", "true") }
)

/*
SendInvoice send invoices. Available Options:
	- OptProviderData(data string)
	- OptPhotoURL(u string)
	- OptPhotoSize(size int)
	- OptPhotoWidth(width int)
	- OptPhotoHeight(height int)
	- OptNeedName
	- OptNeedPhoneNumber
	- OptNeedEmail
	- OptNeedShippingAddress
	- OptSendPhoneNumberToProvider
	- OptSendEmailToProvider
	- OptIsFlexible
	- OptDisableNotification
	- OptReplyToMessageID(id int)
	- OptInlineKeyboardMarkup(markup *InlineKeyboardMarkup)
*/
func (c *Client) SendInvoice(chatID, payload, providerToken string, invoice *Invoice, prices []LabeledPrice, opts ...sendOption) (*Message, error) {
	req := url.Values{}
	req.Set("chat_id", chatID)
	req.Set("title", invoice.Title)
	req.Set("description", invoice.Description)
	req.Set("payload", payload)
	req.Set("provider_token", providerToken)
	req.Set("start_parameter", invoice.StartParameter)
	req.Set("currency", invoice.Currency)
	pr, _ := json.Marshal(prices)
	req.Set("prices", string(pr))
	for _, opt := range opts {
		opt(req)
	}
	msg := &Message{}
	err := c.doRequest("sendInvoice", req, msg)
	return msg, err
}

// ShippingOption represents one shipping option
type ShippingOption struct {
	ID     string         `json:"id"`
	Title  string         `json:"title"`
	Prices []LabeledPrice `json:"prices"`
}

// AnswerShippingQuery options
var (
	OptShippingOptions = func(options []ShippingOption) sendOption {
		return func(v url.Values) {
			op, _ := json.Marshal(options)
			v.Set("shipping_options", string(op))
		}
	}
	OptErrorMessage = func(msg string) sendOption {
		return func(v url.Values) {
			v.Set("error_message", msg)
		}
	}
)

/*
AnswerShippingQuery reply to shipping queries. Available options:
	- OptShippingOptions(options []ShippingOption)
	- OptErrorMessage(msg string)
*/
func (c *Client) AnswerShippingQuery(shippingQueryID string, ok bool, opts ...sendOption) error {
	req := url.Values{}
	req.Set("shipping_query_id", shippingQueryID)
	req.Set("ok", fmt.Sprint(ok))
	for _, opt := range opts {
		opt(req)
	}
	var answered bool
	return c.doRequest("answerShippingQuery", req, &answered)
}

/*
AnswerPreCheckoutQuery respond to pre-checkout queries. Available options:
	- OptErrorMessage(msg string)
*/
func (c *Client) AnswerPreCheckoutQuery(preCheckoutQueryID string, ok bool, opts ...sendOption) error {
	req := url.Values{}
	req.Set("pre_checkout_query_id", preCheckoutQueryID)
	req.Set("ok", fmt.Sprint(ok))
	for _, opt := range opts {
		opt(req)
	}
	var answered bool
	return c.doRequest("answerPreCheckoutQuery", req, &answered)
}

// PassportElementError represents an error in the Telegram Passport element
type PassportElementError interface {
	passportElementError()
}

var (
	_ PassportElementError = PassportElementErrorDataField{}
	_ PassportElementError = PassportElementErrorFrontSide{}
	_ PassportElementError = PassportElementErrorReverseSide{}
	_ PassportElementError = PassportElementErrorSelfie{}
	_ PassportElementError = PassportElementErrorFile{}
	_ PassportElementError = PassportElementErrorFiles{}
)

// PassportElementErrorDataField represents an issue in one of the data fields that was provided by the user
type PassportElementErrorDataField struct {
	Source    string `json:"source"`
	Type      string `json:"type"`
	FieldName string `json:"field_name"`
	DataHash  string `json:"data_hash"`
	Message   string `json:"message"`
}

func (PassportElementErrorDataField) passportElementError() {}

// PassportElementErrorFrontSide represents an issue with the front side of a document
type PassportElementErrorFrontSide struct {
	Source   string `json:"source"`
	Type     string `json:"type"`
	FileHash string `json:"file_hash"`
	Message  string `json:"message"`
}

func (PassportElementErrorFrontSide) passportElementError() {}

// PassportElementErrorReverseSide represents an issue with the reverse side of a document
type PassportElementErrorReverseSide struct {
	Source   string `json:"source"`
	Type     string `json:"type"`
	FileHash string `json:"file_hash"`
	Message  string `json:"message"`
}

func (PassportElementErrorReverseSide) passportElementError() {}

// PassportElementErrorSelfie represents an issue with the selfie with a document
type PassportElementErrorSelfie struct {
	Source   string `json:"source"`
	Type     string `json:"type"`
	FileHash string `json:"file_hash"`
	Message  string `json:"message"`
}

func (PassportElementErrorSelfie) passportElementError() {}

// PassportElementErrorFile represents an issue with a document scan
type PassportElementErrorFile struct {
	Source   string `json:"source"`
	Type     string `json:"type"`
	FileHash string `json:"file_hash"`
	Message  string `json:"message"`
}

func (PassportElementErrorFile) passportElementError() {}

// PassportElementErrorFiles represents an issue with a list of scans
type PassportElementErrorFiles struct {
	Source     string   `json:"source"`
	Type       string   `json:"type"`
	FileHashes []string `json:"file_hashes"`
	Message    string   `json:"message"`
}

func (PassportElementErrorFiles) passportElementError() {}

/*
SetPassportDataErrors informs a user that some of the Telegram Passport elements they provided contains errors
*/
func (c *Client) SetPassportDataErrors(userID int, errors []PassportElementError) error {
	req := url.Values{}
	req.Set("user_id", fmt.Sprint(userID))
	errs, _ := json.Marshal(errors)
	req.Set("errors", string(errs))
	var set bool
	return c.doRequest("setPassportDataErrors", req, &set)
}

/*
SendGame send a game. Available options:
	- OptDisableNotification
	- OptReplyToMessageID(id int)
	- OptInlineKeyboardMarkup(markup *InlineKeyboardMarkup)
*/
func (c *Client) SendGame(chatID, gameShortName string, opts ...sendOption) (*Message, error) {
	req := url.Values{}
	req.Set("chat_id", chatID)
	req.Set("game_short_name", gameShortName)
	for _, opt := range opts {
		opt(req)
	}
	msg := &Message{}
	err := c.doRequest("sendGame", req, msg)
	return msg, err
}

// SetGameScore options
var (
	OptForce = func(v url.Values) {
		v.Set("force", "true")
	}
	OptDisableEditMessage = func(v url.Values) {
		v.Set("disable_edit_message", "true")
	}
)

/*
SetGameScore set the score of the specified user in a game. Available options:
	- OptForce
	- OptDisableEditMessage
*/
func (c *Client) SetGameScore(chatID string, messageID, userID, score int, opts ...sendOption) (*Message, error) {
	req := url.Values{}
	req.Set("chat_id", chatID)
	req.Set("message_id", fmt.Sprint(messageID))
	req.Set("user_id", fmt.Sprint(userID))
	req.Set("score", fmt.Sprint(score))
	for _, opt := range opts {
		opt(req)
	}
	msg := &Message{}
	err := c.doRequest("setGameScore", req, msg)
	return msg, err
}

/*
SetInlineGameScore set the score of the specified user in a game (for inline messages). Available options:
	- OptForce
	- OptDisableEditMessage
*/
func (c *Client) SetInlineGameScore(inlineMessageID string, userID, score int, opts ...sendOption) error {
	req := url.Values{}
	req.Set("inline_message_id", inlineMessageID)
	req.Set("user_id", fmt.Sprint(userID))
	req.Set("score", fmt.Sprint(score))
	for _, opt := range opts {
		opt(req)
	}
	var set bool
	return c.doRequest("setGameScore", req, &set)
}

// GameHighScore represents one row of the high scores table for a game
type GameHighScore struct {
	Position int  `json:"position"`
	User     User `json:"user"`
	Score    int  `json:"score"`
}

/*
GetGameHighScores get data for high score tables
*/
func (c *Client) GetGameHighScores(chatID string, messageID, userID int) ([]*GameHighScore, error) {
	req := url.Values{}
	req.Set("chat_id", chatID)
	req.Set("message_id", fmt.Sprint(messageID))
	req.Set("user_id", fmt.Sprint(userID))
	var scores []*GameHighScore
	err := c.doRequest("getGameHighScores", req, &scores)
	return scores, err
}

/*
GetInlineGameHighScores get data for high score tables
*/
func (c *Client) GetInlineGameHighScores(inlineMessageID string, userID int) ([]*GameHighScore, error) {
	req := url.Values{}
	req.Set("inline_message_id", inlineMessageID)
	req.Set("user_id", fmt.Sprint(userID))
	var scores []*GameHighScore
	err := c.doRequest("getGameHighScores", req, &scores)
	return scores, err
}

/*
SendPoll sends native telegram poll. Available Options:
	- OptDisableNotification
	- OptReplyToMessageID(id int)
	- OptInlineKeyboardMarkup(markup *InlineKeyboardMarkup)
	- OptReplyKeyboardMarkup(markup *ReplyKeyboardMarkup)
	- OptReplyKeyboardRemove
	- OptReplyKeyboardRemoveSelective
	- OptForceReply
	- OptForceReplySelective
*/
func (c *Client) SendPoll(chatID string, question string, options []string, opts ...sendOption) (*Message, error) {
	req := url.Values{}
	req.Set("chat_id", chatID)
	req.Set("question", question)
	marshalledOptions, _ := json.Marshal(options)
	req.Set("options", string(marshalledOptions))
	for _, opt := range opts {
		opt(req)
	}
	msg := &Message{}
	err := c.doRequest("sendPoll", req, msg)
	return msg, err
}

/*
StopPoll stops poll. Available Options:
	- OptInlineKeyboardMarkup(markup *InlineKeyboardMarkup)
*/
func (c *Client) StopPoll(chatID string, messageID string, opts ...sendOption) (*Poll, error) {
	req := url.Values{}
	req.Set("chat_id", chatID)
	req.Set("message_id", messageID)
	for _, opt := range opts {
		opt(req)
	}
	poll := &Poll{}
	err := c.doRequest("stopPoll", req, poll)
	return poll, err
}
