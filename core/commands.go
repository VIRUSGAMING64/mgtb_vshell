package core

import (
	"archive/zip"
	"fmt"
	"govshell/core/downloader"
	"govshell/core/emoji"
	"govshell/core/ffmpeg"
	"net/http"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Response(text string, bot *tgbot.BotAPI, msg *tgbot.Message) error {
	if bot != nil {
		res := tgbot.NewMessage(msg.Chat.ID, text)
		_, err := bot.Send(res)
		return err
	}
	fmt.Println(text)
	return nil
}

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
	bld := strings.Builder{}
	bld.WriteString(text)
	for _, a := range dir {
		_, e := bld.WriteString(a + "\n")
		if e != nil {
			return e
		}
	}
	for _, a := range fils {
		_, e := bld.WriteString(a + "\n")
		if e != nil {
			return e
		}
	}
	text = bld.String()
	err = Response(text, bot, msg)
	return err
}

func mkdir(bot *tgbot.BotAPI, msg *tgbot.Message) error {
	usr := GetUserFromMsg(msg)
	folder := msg.CommandArguments()
	err := usr.Mkdir(folder)
	if err != nil {
		return err
	}
	text := "Directory: " + folder + " created..."
	err = Response(text, bot, msg)
	return err
}

func cd(bot *tgbot.BotAPI, msg *tgbot.Message) error {
	usr := GetUserFromMsg(msg)
	arg := msg.CommandArguments()
	if arg == "" {
		usr.Path = usr.Root
		text := "Directory changed to: [project]/"
		return Response(text, bot, msg)
	}
	err := usr.Chdir(arg)
	if err != nil {
		return err
	}
	text := "Directory changed to: " + arg
	return Response(text, bot, msg)
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
	text := "Removed: " + dir.Name()
	return Response(text, bot, msg)
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

	var m tgbot.Message
	var e error
	if bot != nil {
		m, e = bot.Send(tgbot.NewMessage(msg.Chat.ID, "Compressing..."))
		if e != nil {
			return e
		}
	} else {
		Response("Compressing...", nil, nil)
	}

	err := ffmpeg.Compress(pth, false, func(part, total int64) {
		text := fmt.Sprintf("%.2f%%", float64(part)/float64(total)*100)
		if text == m.Text {
			return
		}
		if bot != nil {
			conf := tgbot.NewEditMessageText(m.Chat.ID, m.MessageID, text)
			m, e = bot.Send(conf)
			if e != nil {
				fmt.Print(e)
			}
			return
		}
		Response(text, bot, msg)

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

func donwl(bot *tgbot.BotAPI, msg *tgbot.Message) error {

	usr := GetUserFromMsg(msg)
	url := msg.CommandArguments()
	dl := downloader.NewDownloader(url)
	if dl.Filename == "" {
		dl.Filename = time.Now().String()
	}
	dl.Filename = path.Join(usr.Path, dl.Filename)
	Response("Downloading", bot, msg)
	dl.StartParallel()
	cli := http.Client{}
	sz, e := downloader.GetLength(&cli, url)
	if e != nil {
		return e
	}
	usr.TDownloaded += sz
	Response("Downloaded", bot, msg)
	return nil
}

func stats(bot *tgbot.BotAPI, msg *tgbot.Message) error {
	mem := new(runtime.MemStats)
	ncpu := runtime.NumCPU()
	ngor := runtime.NumGoroutine()
	runtime.ReadMemStats(mem)
	alloc := mem.Sys / 1024
	mb := false
	if alloc/1024 > 0 {
		alloc /= 1024
		mb = true
	}
	s := "---------BOT STATUS---------\n"
	bld := strings.Builder{}
	bld.WriteString(s)
	bld.WriteString(fmt.Sprint("Alloc memory ", alloc))
	if mb {
		bld.WriteString("MB\n")
	} else {
		bld.WriteString("KB\n")
	}
	bld.WriteString(fmt.Sprintln("CPU Cores:", ncpu))
	bld.WriteString(fmt.Sprintln("Rutines:", ngor))
	text := bld.String()
	return Response(text, bot, msg)
}
