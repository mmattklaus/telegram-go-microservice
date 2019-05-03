package main

import (
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"go-telebot/commands"
	tb "gopkg.in/tucnak/telebot.v2"
	l "log"
	"os"
	"time"
)

var (
	BotName = os.Getenv("TB_NAME")
	BotUri = os.Getenv("TB_URI")
	BotToken = os.Getenv("TB_TOKEN")
)

func main() {
	fmt.Println("Initializing go telebot...")

	log := l.New(os.Stdout, "TB: ", l.LstdFlags | l.Lshortfile)
	bot, err := tb.NewBot(tb.Settings{
		Token:  BotToken,
		URL: BotUri,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatalf("Unable to connect to bot server: %v", err)
	}

	bh := commands.HandleCMD(bot, log, &BotName)

	bot.Handle("/hello", func(m *tb.Message) {
		go bh.HandleGreeting(m)
	})

	bot.Handle("/record", func(m *tb.Message) {
		go bh.HandleRecording(m)
	})

	bot.Handle("/pic", func(m *tb.Message) {
		go bh.HandleSnap(m)
	})

	bot.Start()

}
