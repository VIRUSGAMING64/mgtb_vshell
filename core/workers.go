package core

import (
	"sync"
	"time"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var comp_queue []*tgbot.BotAPI
var msg_queue []*tgbot.Message

func CompressingWorker() {
	for {
		time.Sleep(time.Second * 3)
		if len(msg_queue) == 0 {
			continue
		}
		bot := comp_queue[0]
		msg := msg_queue[0]
		err := x265(bot, msg)
		if err != nil {
			bot.Send(tgbot.NewMessage(msg.Chat.ID, err.Error()))
		}
		comp_queue = comp_queue[1:]
		msg_queue = msg_queue[1:]
	}
}

var saveMTX = sync.Mutex{}

func SaveWorker() {
	for {
		time.Sleep(1 * time.Minute)
		saveMTX.Lock()
		e := SaveUsers()
		if e != nil {
			panic(e)
		}
		saveMTX.Unlock()
	}
}
