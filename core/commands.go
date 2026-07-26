package core

import (
	"archive/zip"
	"fmt"
	"govshell/core/emoji"
	"govshell/core/ffmpeg"
	"os"
	"path"

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

func rm(bot *tgbot.BotAPI, msg *tgbot.Message) error {
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
	err := usr.Rm(dir.Name())
	if err != nil {
		return err
	}
	res := tgbot.NewMessage(msg.Chat.ID, "Removed: "+dir.Name())
	_, e := bot.Send(res)
	return e
}

func x265(bot *tgbot.BotAPI, msg *tgbot.Message) error {
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
	pth := path.Join(usr.Path, dir.Name())
	m, e := bot.Send(tgbot.NewMessage(msg.Chat.ID, "Compressing..."))
	if e != nil {
		return e
	}
	err := ffmpeg.Compress(pth, false, func(part, total int64) {
		fmt.Println(total, part)
		text := fmt.Sprintf("%.2f%%", float64(part)/float64(total)*100)
		if text == m.Text {
			return
		}
		conf := tgbot.NewEditMessageText(m.Chat.ID, m.MessageID, text)
		m, e = bot.Send(conf)
		fmt.Println(e)
	})
	return err
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

	fs := os.DirFS(path.Join(usr.Path, dir.Name()))
	f, e := os.OpenFile(dir.Name()+".zip", os.O_CREATE|os.O_RDWR, 0644)
	if e != nil {
		return e
	}
	defer f.Close()

	res := tgbot.NewMessage(msg.Chat.ID, "Compressing")
	bot.Send(res)
	zfile := zip.NewWriter(f)
	defer zfile.Close()

	zfile.AddFS(fs)
	res = tgbot.NewMessage(msg.Chat.ID, "Compressed")
	bot.Send(res)
	return nil
}
