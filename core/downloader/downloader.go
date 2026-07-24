package downloader

import (
	"errors"
	"net/http"
	"os"
	"time"
)

type Worker struct {
	L   int
	R   int
	Url string
}

type Downlaoder struct {
	Url      string
	threads  byte
	mthreads byte
	timeout  time.Duration
	Cli      *http.Client
	workers  []Worker
	Filename string
}

var chunk_size int64 = 512 * 1024

func (a *Downlaoder) StartSec() error {
	/*
		Start a sequential download of url
	*/
	if a.Filename == "" {
		return errors.New("Empty filename")
	}
	ContentSize, e := GetLength(a.Cli, a.Url)
	if ContentSize == -1 {
		return e
	}
	writer, e := os.OpenFile(a.Filename, os.O_CREATE|os.O_RDWR, 0644)
	if e != nil {
		return e
	}
	defer writer.Close()
	i := int64(0)
	sz := 0

	for ; i+chunk_size <= ContentSize; i += chunk_size {
		data, err := DlWoker(a.Url, i, i+chunk_size-1)
		sz += len(data)
		if err != nil {
			return err
		}
		_, err = writer.Write(data)
		if err != nil {
			return err
		}
	}
	data, err := DlWoker(a.Url, i, ContentSize)
	sz += len(data)
	if err != nil {
		return err
	}
	_, err = writer.Write(data)
	return err
}

func NewDownloader(url string) *Downlaoder {

	cli := http.Client{
		Timeout: time.Second * 3,
	}
	dl := Downlaoder{Url: url, timeout: 3 * time.Second, Cli: &cli}

	return &dl
}
