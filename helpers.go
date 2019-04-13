package tbot

// Buttons construct ReplyKeyboardMarkup from strings
func Buttons(buttons [][]string) *ReplyKeyboardMarkup {
	keyboard := make([][]KeyboardButton, len(buttons))
	for i := range buttons {
		keyboard[i] = make([]KeyboardButton, len(buttons[i]))
		for j := range buttons[i] {
			keyboard[i][j] = KeyboardButton{Text: buttons[i][j]}
		}
	}
	return &ReplyKeyboardMarkup{Keyboard: keyboard}
}
