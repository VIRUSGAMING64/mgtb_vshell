package main

import (
	"context"
	"fmt"
	"govshell/core"
	"govshell/web"
	"os"
	"os/signal"
	"syscall"
	"time"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	go web.Start()
	go core.Aliver()
	go core.InitConfig()
	time.Sleep(60 * time.Second)
	BOT_TOKEN, exsits := os.LookupEnv("BOT_TOKEN")

	if exsits == false {
		fmt.Println("bot token not present")
		os.Exit(1)
	}

	go func() {
		ctx, stop := signal.NotifyContext(
			context.Background(), syscall.SIGINT)
		defer stop()
		<-ctx.Done()
		core.SaveUsers()
		fmt.Println("\nProcess ended: ", context.Cause(ctx))
		os.Exit(int(syscall.SIGINT))
	}()

	bot, err := tgbot.NewBotAPI(BOT_TOKEN)
	core.Bot = bot
	go core.MainWorker(bot)
	if err != nil {
		core.Offline = true
		fmt.Println("Starting offline")
		core.Cmd()
		fmt.Println(err)
		return
	}

	fmt.Println("Bot started")
	os.MkdirAll("roots", 0644)
	updatesConf := tgbot.NewUpdate(0)
	updatesConf.Timeout = 60
	updatesChan := bot.GetUpdatesChan(updatesConf)

	for update := range updatesChan {
		core.Updates <- update
	}
}
