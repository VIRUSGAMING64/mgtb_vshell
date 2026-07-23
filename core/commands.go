package core

import (
	"govshell/core/emoji"
	"govshell/core/ffmpeg"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ls(bot *tgbot.BotAPI, msg *tgbot.Message) error {
	usr := GetUserFromMsg(msg)
	dirs, err := usr.ListDir()
	if err != nil {
		return err
	}
	var dir []string
	var fils []string
	for _, a := range dirs {
		if a.IsDir() {
			dir = append(dir, emoji.FILE_FOLDER+" "+a.Name())
		} else {
			fils = append(
				fils,
				emoji.PAGE_FACING_UP+" "+a.Name(),
			)
		}
	}

	text := usr.Path + "\n"
	for _, a := range dir {
		text += a + "\n"
	}
	for _, a := range fils {
		text += a + "\n"
	}
	res := tgbot.NewMessage(msg.Chat.ID, text)
	_, err = bot.Send(res)
	return err
}

func mkdir(bot *tgbot.BotAPI, msg *tgbot.Message) error {
	usr := GetUserFromMsg(msg)
	folder := msg.CommandArguments()
	err := usr.Mkdir(folder)
	if err != nil {
		return err
	}
	res := tgbot.NewMessage(msg.Chat.ID, "Directory: "+folder+" created...")
	_, err = bot.Send(res)
	return err
}

// TODO
func x265(bot *tgbot.BotAPI, msg *tgbot.Message) error {
	return nil
}

// TODO
func upload(bot *tgbot.BotAPI, msg *tgbot.Message) error {
	return nil
}

func zips(bot *tgbot.BotAPI, msg *tgbot.Message) error {
	usr := GetUserFromMsg(msg)
	arg := msg.CommandArguments()
	idx := ffmpeg.String2Int64(arg)
	dir, e1 := GetFromIndex(idx, usr)
	di2, e2 := GetFromName(arg, usr)
	if dir == nil {
		dir = di2
		e1 = e2
	}

	if e1 != nil {
		return e1
	}

	//f, _ := os.Open("tmp")
	//zfile := zip.NewWriter(f)

	//TODO Aqui tengo que hacerle zip a ese directory

	return nil
}
