package commands

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	e "go-telebot/ems"
	f "go-telebot/functions"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

/**
* @todo Map cmds to a list of possible alternatives
* @todo Save request to database in Middleware <Monitor>
*/

type BotHandler struct {
	ID   *string
	Bot  *tb.Bot
	log  *log.Logger
	db   *sqlx.DB
}

func HandleCMD(tb *tb.Bot, l *log.Logger, id *string, db *sqlx.DB) *BotHandler {
	return &BotHandler{
		ID:  id,
		Bot: tb,
		log: l,
		db:  db,
	}
}

func (bh *BotHandler) SetupRoutes() {
	bh.Bot.Handle("/start", bh.Monitor(bh.HandleWelcome))

	bh.Bot.Handle("/hello", bh.Monitor(bh.HandleGreeting))

	bh.Bot.Handle("/record", bh.Monitor(bh.HandleRecording))

	bh.Bot.Handle("/pic", bh.Monitor(bh.HandleSnap))

	bh.Bot.Handle("/ip", bh.Monitor(bh.HandleIP))

	// General catchers
	bh.Bot.Handle(tb.OnText, func(m *tb.Message) {
		// all the text messages that weren't
		// captured by existing handlers
	})

	bh.Bot.Handle(tb.OnPhoto, func(m *tb.Message) {
		// photos only
	})

	bh.Bot.Handle(tb.OnChannelPost, func (m *tb.Message) {
		// channel posts only
	})

	bh.Bot.Handle(tb.OnQuery, func (q *tb.Query) {
		// incoming inline queries
	})
}

func (bh *BotHandler) HandleGreeting(m *tb.Message) {
	_, err := bh.Bot.Reply(m, fmt.Sprintf("Hello, %s %s", m.Sender.FirstName, e.Ems("wink", "smile", "cloud")))
	if err != nil {
		bh.log.Fatalln(err)
	}
}
func (bh *BotHandler) HandleSnap(m *tb.Message) {
	filename := fmt.Sprintf("%d.jpg", time.Now().UnixNano())
	bh.log.Printf("snapping image")
	_, _ = bh.Bot.Reply(m, fmt.Sprintf("%s...", e.Ems("camera")))
	// snap photo with commandline
	err :=  f.Snap(filename)
	// load photo
	photo := &tb.Photo{File: tb.FromDisk(filename)}
	// Sender photo as a reply
	_, err = bh.Bot.Reply(m, photo)
	if err != nil {
		bh.log.Fatalln(err)
	}
	// Delete image after sending
	err = os.Remove(filename)
	if err != nil {
		bh.log.Fatalln(err)
	}
}

func (bh *BotHandler) HandleRecording(m *tb.Message) {
	wd, err := os.Getwd()
	payload := strings.Split(m.Payload, " ")
	duration, err := strconv.Atoi(payload[0])
	if err != nil {
		bh.log.Println("Duration not specified")
		duration = 3
	}
	filename := fmt.Sprintf("%d.wav", time.Now().UnixNano())

	go func() {
		err := f.Record(filename, duration)
		if err != nil {
			defer bh.log.Printf("unable to record audio: %v", err)
			_, err = bh.Bot.Reply(m, "Sorry, I couldn't record the audio.")
			if err != nil {
				bh.log.Fatalf("Error message not sent to user: %v", m.Sender)
			}
		}
		bh.log.Println("audio recorded. Sending")
		path := fmt.Sprintf("%s/%s", wd, filename)
		audio := &tb.Audio{File: tb.FromDisk(path)}
		_, err = bh.Bot.Send(m.Sender, audio)
		if err != nil {
			bh.log.Printf("unable to send audio to: %v", m.Sender.FirstName)
		}
		err = os.Remove(path)
	}()
	_, err = bh.Bot.Send(m.Sender, fmt.Sprintf("%s...", e.Ems("microphone")), &tb.SendOptions{
		ReplyTo: m,
	})
}

func (bh *BotHandler) HandleIP(m *tb.Message) {
	ip, err := f.Ip()
	if len(ip) > 0 {
		_, err := bh.Bot.Reply(m, ip)
		if err != nil {
			log.Fatalf("err sending ip addr: %v", err)
		}
	}
	if err != nil {
		log.Fatalf("error in handle cmd: %v", err)
	}
}

func (bh *BotHandler) HandleWelcome(m *tb.Message) {
	_, err := bh.Bot.Reply(m, fmt.Sprintf("I personnal welcome you! %s %s", e.Ems("wink", "heart"), *bh.ID))
	if err != nil {
		log.Fatalf("unable to send welcome message: %v", err)
	}
}

/*
* Middleware to handle requests
*/
func (bh *BotHandler) Monitor(next func(m *tb.Message)) func(m *tb.Message) {
	return func(m *tb.Message) {
		startTime := time.Now()
		user := m.Sender.FirstName
		go next(m)
		defer bh.log.Printf("Request made by {%s} in %s", user, time.Since(startTime))
	}
}