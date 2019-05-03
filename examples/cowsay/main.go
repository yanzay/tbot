package main

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/yanzay/tbot/v2"
)

func main() {
	token := os.Getenv("TELEGRAM_TOKEN")
	bot := tbot.New(token)
	c := bot.Client()
	bot.HandleMessage("cowsay .+", func(m *tbot.Message) {
		text := strings.TrimPrefix(m.Text, "cowsay ")
		cow := fmt.Sprintf("```\n%s\n```", cowsay(text))
		c.SendMessage(m.Chat.ID, cow, tbot.OptParseModeMarkdown)
	})
	bot.Start()
}

func cowsay(text string) string {
	lineLen := utf8.RuneCountInString(text) + 2
	topLine := fmt.Sprintf(" %s ", strings.Repeat("_", lineLen))
	textLine := fmt.Sprintf("< %s >", text)
	bottomLine := fmt.Sprintf(" %s ", strings.Repeat("-", lineLen))
	cow := `
        \   ^__^
         \  (oo)\_______
            (__)\       )\/\
               ||----w |
               ||     ||
	`
	resp := fmt.Sprintf("%s\n%s\n%s%s", topLine, textLine, bottomLine, cow)
	return resp
}
