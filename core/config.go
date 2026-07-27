package core

import (
	"fmt"
	"govshell/core/ffmpeg"
	"io"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	ADMIN = 1
	BAN   = 2
)

var ROOT = "roots/"
var SAVE_USERS = "users.json"
var AdminsPth = "admins.conf"
var Admins []int64
var Updates chan tgbotapi.Update = make(chan tgbotapi.Update)

func InitConfig() {
	f, e := os.OpenFile(AdminsPth, os.O_RDONLY, 0644)
	if e != nil {
		fmt.Println(e.Error())
	} else {
		data, err := io.ReadAll(f)
		if err != nil {
			fmt.Println(err)
		} else {
			ts := string(data)
			for _, ele := range strings.Split(ts, "\n") {
				id := ffmpeg.String2Int64(ele)
				Admins = append(Admins, id)
			}
		}
	}
	defer f.Close()
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
