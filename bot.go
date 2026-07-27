package main

import (
	"context"
	"fmt"
	"govshell/core"
	"os"
	"os/signal"
	"syscall"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	core.InitConfig()
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

	go core.MainWorker(bot)
	if err != nil {
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
