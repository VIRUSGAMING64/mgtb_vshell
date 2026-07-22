package core

import (
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetUserFromMsg(msg *tgbot.Message) *User {
	return Users[msg.From.ID]
}
