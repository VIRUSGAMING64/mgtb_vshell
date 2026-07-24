package downloader

import (
	"errors"
	"govshell/core/ffmpeg"
	"io"
	"net/http"
	"time"
)

func DlWoker(url string, l, r int64) ([]byte, error) {
	cli := http.Client{
		Timeout: time.Second * 5,
	}
	res, e := cli.Head(url)
	if e != nil {
		return nil, e
	}
	slen := res.Header.Get("Content-Length")
	len := ffmpeg.String2Int64(slen)
	if len == -1 {
		return nil, errors.New("Not detected content size")
	}
	ranges := res.Header.Get("Accept-Ranges")
	if ranges != "bytes" {
		return nil, errors.New("Not seekable file")
	}
	req, e := http.NewRequest("GET", url, nil)
	if e != nil {
		return nil, e
	}
	req.Header.Add("Range", "bytes="+ffmpeg.Int64toString(l)+"-"+ffmpeg.Int64toString(r))
	res, e = cli.Do(req)
	if e != nil {
		return nil, e
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}
