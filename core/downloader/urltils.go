package downloader

import (
	"errors"
	"govshell/core/ffmpeg"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func GetLength(cli *http.Client, url string) (int64, error) {
	res, e := cli.Head(url)
	if e != nil {
		return -1, e
	}
	slen := res.Header.Get("Content-Length")
	len := ffmpeg.String2Int64(slen)
	if len == -1 {
		return -1, errors.New("Not detected content size")
	}
	return len, nil
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
