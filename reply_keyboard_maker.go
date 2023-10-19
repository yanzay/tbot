package tbot

type replyRow struct {
	buttons []KeyboardButton
}

// Short version for adding a new button to a replyRow,
// just accepts a text and ignore other properties of a button,
// which is the common case.
func (r *replyRow) AddButton(text string) {
	r.buttons = append(r.buttons, KeyboardButton{
		Text: text,
	})
}

// Full versions of adding a new button
func (r *replyRow) AddButtonFull(button KeyboardButton) {
	r.buttons = append(r.buttons, button)
}

type replyKeyboardMarker struct {
	replyRows []*replyRow
	resize    bool
	oneTime   bool
	selective bool
}

// NewReplyKeyboardMaker makes a replyKeyboardMarker (builder design pattern)
// which then can be used for making reply keyboards,
//
// It has some methods for adding new rows and you can call Build method at the end and use returned Keyboard.
func NewReplyKeyboardMaker() *replyKeyboardMarker {
	return &replyKeyboardMarker{}
}

// Add a new row to a Reply keyboard and returns a pointer to it,
// then you can use this pointer to add new buttons using `.AddButton` or `.AddButtonFull`
func (km *replyKeyboardMarker) AddRow() *replyRow {
	newReplyRow := new(replyRow)
	km.replyRows = append(km.replyRows, newReplyRow)
	return newReplyRow
}

func (km *replyKeyboardMarker) Build() *ReplyKeyboardMarkup {
	ikm := &ReplyKeyboardMarkup{
		Keyboard:        make([][]KeyboardButton, len(km.replyRows)),
		ResizeKeyboard:  km.resize,
		OneTimeKeyboard: km.oneTime,
		Selective:       km.selective,
	}
	for i, replyRow := range km.replyRows {
		ikm.Keyboard[i] = replyRow.buttons
	}
	return ikm
}

func (km *replyKeyboardMarker) SetResize(resize bool) {
	km.resize = resize
}

func (km *replyKeyboardMarker) SetSelective(selective bool) {
	km.selective = selective
}

func (km *replyKeyboardMarker) SetOneTime(oneTime bool) {
	km.oneTime = oneTime
}
