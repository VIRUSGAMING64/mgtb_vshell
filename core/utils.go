package core

import (
	"errors"
	"net/url"
	"os"
	"strings"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetUserFromMsg(msg *tgbot.Message) *User {
	return Users[msg.From.ID]
}

func GetFromIndex(idx int64, usr *User) (os.DirEntry, error) {
	dirs, err := usr.ListDir()
	if err != nil {
		return nil, err
	}
	if idx < 0 || idx >= int64(len(dirs)) {
		return nil, errors.New("Invalid index")
	}
	return dirs[idx], nil
}

func GetFromName(arg string, usr *User) (os.DirEntry, error) {
	dirs, e := usr.ListDir()
	if e != nil {
		return nil, e
	}
	for _, a := range dirs {
		if a.Name() == arg {
			return a, nil
		}
	}
	return nil, errors.New("Not found entry")
}

func UriEncode2String(s string) string {
	s = strings.ReplaceAll(s, "%20", " ")
	s = strings.ReplaceAll(s, "%28", "(")
	s = strings.ReplaceAll(s, "%29", ")")
	return s
}

func GetNameFromUrl(uri url.URL) string {
	pths := strings.Split(uri.Path, "/")
	if len(pths) == 0 {
		return ""
	}
	return UriEncode2String(pths[len(pths)-1])
}

func JoinFiles(pth string, names []string) error {
	e := os.RemoveAll(pth)
	if e != nil {
		return e
	}
	writer, err := os.OpenFile(pth, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer writer.Close()
	for _, name := range names {
		f, e := os.OpenFile(name, os.O_RDONLY, 0644)
		if e != nil {
			return e
		}
		defer f.Close()
		for {
			var buff []byte = make([]byte, 512*1024)
			cnt, e := f.Read(buff)
			buff = buff[:cnt]
			if e != nil {
				if e.Error() == "EOF" {
					break
				}
				return e
			}
			_, e = writer.Write(buff)
			if e != nil {
				return e
			}
		}
		e = os.RemoveAll(name)
		if e != nil {
			return e
		}
	}
	return nil
}
