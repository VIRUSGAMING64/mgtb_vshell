package downloader

import (
	"errors"
	"govshell/core/ffmpeg"
	"net/http"
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
