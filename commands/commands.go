package commands

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type BotHandler struct {
	ID *string
	bot *tb.Bot
	log *log.Logger
}

func HandleCMD(tb *tb.Bot, l *log.Logger, id *string) *BotHandler {
	return &BotHandler{
		ID: id,
		bot: tb,
		log: l,
	}
}

func (bh *BotHandler) HandleGreeting(m *tb.Message) {
	fmt.Printf("Firstname: %+v  Lastname: %s  ID: %v  Username: %v LangCode: %v\n", m.Sender.FirstName, m.Sender.LastName, m.Sender.ID, m.Sender.Username, m.Sender.LanguageCode)
	_, err := bh.bot.Reply(m, fmt.Sprintf("hello, %s", m.Sender.FirstName))
	if err != nil {
		bh.log.Fatalln(err)
	}
}
func (bh *BotHandler) HandleSnap(m *tb.Message) {
		filename := fmt.Sprintf("%d.jpg", time.Now().UnixNano())
		bh.log.Printf("snapping image")
		// snap photo with commandline
		err := snap(filename)
		// load photo
		photo := &tb.Photo{File: tb.FromDisk(filename)}
		// Sender photo as a reply
		_, err = bh.bot.Reply(m, photo)
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
		bh.log.Println("unable to convert duration")
		duration = 3
	}
	filename := fmt.Sprintf("%d.wav", time.Now().UnixNano())

	func() {
		err := record(filename, duration)
		if err != nil {
			defer bh.log.Println("unable to record audio.")
			_, err = bh.bot.Reply(m, "Sorry, I couldn't record the audio.")
			if err != nil {
				bh.log.Fatalf("Error message not sent to user: %v", m.Sender)
			}
		}
		bh.log.Println("audio recorded. Sending")
		path := fmt.Sprintf("%s/%s", wd, filename)
		audio := &tb.Audio{File: tb.FromDisk(path)}
		_, err = bh.bot.Send(m.Sender, audio)
		if err != nil {
			bh.log.Printf("unable to send audio to: %v", m.Sender.FirstName)
		}
	}()
	_, err = bh.bot.Send(m.Sender, "recording...", &tb.SendOptions{
		ReplyTo: m,
	})
}

func record(filename string, duration int) error {
	if duration > 5 * 60 {
		fmt.Println("duration too lengthy")
		duration = 3
	}
	fmt.Println(duration)
	cmd := exec.Command("rec", "-r", "160000", "-c", "1", filename , "trim", "0", string(duration)) //

	env := os.Environ()
	// env = append(env, "AUDIODEV=hw:1,0")
	cmd.Env = env
	fmt.Println("recording audio...")
	return cmd.Run()
}

func snap(filename string) (error) {
	cmd := exec.Command("fswebcam", "--no-banner",  filename)
	return cmd.Run()
}