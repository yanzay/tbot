package main

import (
	"log"
	"os"
	"time"

	"github.com/yanzay/tbot"
)

func main() {
	token := os.Getenv("TELEGRAM_TOKEN")
	// Create new telegram bot server using token
	bot, err := tbot.NewServer(token)
	if err != nil {
		log.Fatal(err)
	}

	//middleware works selected auth options usernames,userids or chatids
	//username list
	//whitelist := []string{"user1", "user2", "user3"}
	//or
	//userid list (you can use this option to grant access users)
	//whitelist := []int{521844285}
	//or
	//chatid list (you can use this option to grant access groups)
	whitelist := []int64{-271651159, -266451058, 512718860}
	bot.AddMiddleware(tbot.NewAuth(whitelist))

	// Handle with HiHandler function
	bot.HandleFunc("/hi", HiHandler)
	err = bot.ListenAndServe()
	log.Fatal(err)
}

func HiHandler(message *tbot.Message) {
	// Handler can reply with several messages
	message.Replyf("Hello, %s!", message.From.FirstName)
	time.Sleep(1 * time.Second)
	message.Reply("What's up?")
}
