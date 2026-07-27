package core

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

const (
	ADMIN = 1
	BAN   = 2
)

var ROOT = "roots/"
var SAVE_USERS = "users.json"
var Admins []int64
var Updates chan tgbotapi.Update = make(chan tgbotapi.Update)

func InitConfig() {
	Admins = append(Admins, 0)
	LoadUsers()
	go SaveWorker()
	go CompressingWorker()
	for _, id := range Admins {
		if Users[id] == nil {
			Users[id] = NewUser(id)
		}
		Users[id].Status |= ADMIN
		Users[id].Root = "."

	}
	SaveUsers()
}
