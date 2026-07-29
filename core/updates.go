package core

import (
	"fmt"
	"sync"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var updatesmtx sync.Mutex = sync.Mutex{}

func updateUsers(msg *tgbot.Message) {
	updatesmtx.Lock()
	defer updatesmtx.Unlock()
	usrid := msg.From.ID
	usr := msg.From.UserName
	if Users[(usrid)] == nil {
		Users[(usrid)] = NewUser((usrid))
	}
	Users[(usrid)].Name = usr
}

func OnCommand(bot *tgbot.BotAPI, msg *tgbot.Message) error {
	cmd := msg.Command()
	var err error = nil
	switch cmd {
	case "zip":
		err = zips(bot, msg)
	case "x265":
		Response(fmt.Sprintf("Queue position: %d", len(msg_queue)), bot, msg)
		msg_queue = append(msg_queue, msg)
	case "stats":
		err = stats(bot, msg)
	case "ls":
		err = ls(bot, msg)
	case "cd":
		err = cd(bot, msg)
	case "mkdir":
		err = mkdir(bot, msg)
	case "upload":
		err = upload(bot, msg)
	case "rm":
		err = rm(bot, msg)
	case "downl":
		err = donwl(bot, msg)
	default:
		fmt.Println("unknow command:", cmd)
	}
	if err != nil {
		fmt.Println(err.Error())
	}
	return nil
}

func OnInlineQuery(bot *tgbot.BotAPI, q *tgbot.InlineQuery) error {
	return nil
}

func OnMessage(bot *tgbot.BotAPI, msg *tgbot.Message) error {
	updateUsers(msg)

	adm := Users[msg.From.ID]
	if adm.Status&BAN != 0 {
		fmt.Println("user banned: ", adm.Id)
		return nil
	}
	if msg.IsCommand() {
		return OnCommand(bot, msg)
	}
	fmt.Println("Normal message:", msg.Text)
	return nil
}

func OnUpdate(bot *tgbot.BotAPI, up tgbot.Update) error {
	if up.Message != nil {
		return OnMessage(bot, up.Message)
	}
	if up.InlineQuery != nil {
		return OnInlineQuery(bot, up.InlineQuery)
	}
	fmt.Println("Unknow update")
	return nil
}
