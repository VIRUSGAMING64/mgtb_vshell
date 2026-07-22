package core

import (
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
		Users[(usrid)].Name = usr
	}
}

func OnCommand(bot *tgbot.BotAPI, msg *tgbot.Message) error {
	cmd := msg.Command()
	out := ""
	switch cmd {
	case "x265":
		x265(bot, msg)
	case "ls":
		ls(bot, msg)
	case "mkdir":
		mkdir(bot, msg)
	case "upload":
		upload(bot, msg)
	}
	bot.Send(tgbot.NewMessage(msg.Chat.ID, out))
	return nil
}

func OnInlineQuery(bot *tgbot.BotAPI, q *tgbot.InlineQuery) error {
	return nil
}

func OnMessage(bot *tgbot.BotAPI, msg *tgbot.Message) error {
	updateUsers(msg)
	if msg.IsCommand() {
		return OnCommand(bot, msg)
	}
	return nil
}

func OnUpdate(bot *tgbot.BotAPI, up tgbot.Update) error {
	if up.Message != nil {
		return OnMessage(bot, up.Message)
	}
	if up.InlineQuery != nil {
		return OnInlineQuery(bot, up.InlineQuery)
	}
	return nil
}
