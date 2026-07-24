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
