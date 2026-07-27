package core

import (
	"fmt"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Cmd() {
	for {
		fmt.Printf("> ")
		var s string
		var buff [4096]byte
		n, e := os.Stdin.Read(buff[:])
		if e != nil {
			fmt.Println(e)
		}
		s = string(buff[:n-1])
		if s == "" {
			continue
		}
		up := tgbotapi.Update{}
		msg := tgbotapi.Message{}
		up.Message = &msg
		up.Message.Text = s
		usr := tgbotapi.User{}
		msg_ent := tgbotapi.MessageEntity{}
		if strings.HasPrefix(s, "/") {
			msg_ent.Type = "bot_command"
		}
		msg.CommandArguments()
		up.Message.Entities = make([]tgbotapi.MessageEntity, 0)
		up.Message.Entities = append(up.Message.Entities, msg_ent)
		up.Message.Entities[0].Length = len(strings.Split(s, " ")[0])
		up.Message.From = &usr
		up.Message.From.UserName = "offline"
		up.Message.From.ID = 0
		Updates <- up
		time.Sleep(1 * time.Second)
	}
}

func Line() {

}
