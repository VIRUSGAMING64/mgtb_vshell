package main

import (
	"fmt"
	"govshell/core"
	"os"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	core.InitConfig()
	BOT_TOKEN, exsits := os.LookupEnv("BOT_TOKEN")
	if exsits == false {
		fmt.Println("bot token not present")
		os.Exit(1)
	}
	bot, err := tgbot.NewBotAPI(BOT_TOKEN)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Bot started")
	core.LoadUsers()
	go core.SaveWorker()
	go core.CompressingWorker()
	os.MkdirAll("roots", 0644)
	updatesConf := tgbot.NewUpdate(0)
	updatesConf.Timeout = 60
	updatesChan := bot.GetUpdatesChan(updatesConf)
	for update := range updatesChan {
		go func() {
			err = core.OnUpdate(bot, update)
			if err != nil {
				fmt.Println(err)
			}
		}()
	}
}
