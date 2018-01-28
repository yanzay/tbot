package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/eeonevision/tbot"
)

func main() {
	bot, err := tbot.NewServer(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}
	bot.HandleFunc("/cowsay {text}", CowHandler)
	err = bot.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func CowHandler(m *tbot.Message) {
	reply := fmt.Sprintf("```\n%s\n```", cowsay(m.Vars["text"]))
	m.Reply(reply, tbot.WithMarkdown)
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
