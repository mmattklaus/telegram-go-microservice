package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
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
	db := sqlx.MustConnect("sqlite3", ":memory:")

	bh := commands.HandleCMD(bot, log, &BotName, db)

	bh.SetupRoutes()

	bh.Bot.Start()

}
