package tbot

type row struct {
	buttons []InlineKeyboardButton
}

// Short version for adding a new button to a row
func (r *row) AddButton(text, data, url string) {
	r.buttons = append(r.buttons, InlineKeyboardButton{
		Text:         text,
		CallbackData: data,
		URL:          url,
	})
}

func (r *row) AddButtonFull(button InlineKeyboardButton) {
	r.buttons = append(r.buttons, button)
}

type keyboardMarker struct {
	rows []*row
}

// NewKeyboardMaker makes a keyboardMarker (builder design pattern)
// which then can be used for making inline keyboards,
//
// It has some methods for adding new rows and you can call Build method at the end and use returned Keyboard.
func NewKeyboardMaker() *keyboardMarker {
	return &keyboardMarker{}
}

// Add a new row to an inline keyboard and returns a pointer to it,
// then you can use this pointer to add new buttons using `.AddButton` or `.AddButtonFull`
func (km *keyboardMarker) AddRow() *row {
	newRow := new(row)
	km.rows = append(km.rows, newRow)
	return newRow
}

func (km *keyboardMarker) Build() *InlineKeyboardMarkup {
	ikm := &InlineKeyboardMarkup{
		InlineKeyboard: make([][]InlineKeyboardButton, len(km.rows)),
	}
	for i, row := range km.rows {
		ikm.InlineKeyboard[i] = row.buttons
	}
	return ikm
}
