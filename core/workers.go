package core

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var msg_queue []*tgbot.Message

func CompressingWorker() {
	for {
		time.Sleep(time.Second * 3)
		if len(msg_queue) == 0 {
			continue
		}
		msg := msg_queue[0]
		err := x265(Bot, msg)
		if err != nil {
			Response(err.Error(), Bot, msg)
		}
		msg_queue = msg_queue[1:]
		if len(msg_queue) != 0 {
			for i := range msg_queue {
				tmsg := msg_queue[i]
				Response(fmt.Sprintf("Queue position of %s changed to: %d", msg.Text, i), Bot, tmsg)
				time.Sleep(200 * time.Millisecond)
			}
		}
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

func MainWorker(bot *tgbot.BotAPI) {
	for update := range Updates {
		go OnUpdate(bot, update)
	}
}

func Aliver() {
	time.Sleep(2 * time.Second)
	for {
		_, e := http.Get("http://localhost/alive")
		if e != nil {
			fmt.Println(e.Error())
		}
		time.Sleep(time.Second * 30)
	}
}
