package downloader

import (
	"errors"
	"fmt"
	"govshell/core/ffmpeg"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

type Worker struct {
	L   int
	R   int
	Url string
}

type Downloader struct {
	Url      string
	threads  byte
	running  byte
	workers  []Worker
	Filename string
}

var dlmutex sync.Mutex
var dlerrors chan error
var chunk_size int64 = 512 * 1024

func (a *Downloader) StartSec() error {
	/*
		Start a sequential download of url
	*/
	cli := http.Client{
		Timeout: time.Second * 3,
	}
	ContentSize, e := GetLength(&cli, a.Url)
	if ContentSize == -1 {
		return e
	}
	return ParallelWorker(a, a.Filename, a.Url, 0, ContentSize)
}

func NewDownloader(uri string) *Downloader {
	dl := Downloader{Url: uri, threads: 8}
	ur, e := url.Parse(uri)
	if e != nil {
		fmt.Println(e)
		panic(e)
	}
	dl.Filename = GetNameFromUrl(*ur)
	return &dl
}

func ParallelWorker(a *Downloader, file, url string, l, r int64) error {
	dlmutex.Lock()
	a.running += 1
	dlmutex.Unlock()
	defer func() {
		dlmutex.Lock()
		a.running -= 1
		dlmutex.Unlock()
	}()
	if file == "" {
		return errors.New("Empty filename")
	}
	e := os.RemoveAll(file)
	if e != nil {
		return e
	}
	writer, e := os.OpenFile(file, os.O_CREATE|os.O_RDWR, 0644)
	if e != nil {
		return e
	}
	defer writer.Close()
	i := l
	sz := int64(0)

	for ; i+chunk_size <= r; i += chunk_size {
		data, err := DlWoker(url, i, i+chunk_size-1)
		sz += int64(len(data))
		if err != nil {
			return err
		}
		_, err = writer.Write(data)
		if err != nil {
			return err
		}
	}
	data, err := DlWoker(url, i, r)
	sz += int64(len(data))
	if err != nil {
		return err
	}
	_, err = writer.Write(data)
	return err
}

func (a *Downloader) StartParallel() error {
	cli := http.Client{}
	if a.Filename == "" {
		return errors.New("Empty filename")
	}
	ContentSize, e := GetLength(&cli, a.Url)
	BytesPrThread := ContentSize / int64(a.threads)
	if ContentSize == -1 {
		return e
	}

	var names []string
	for i := int64(0); i < int64(a.threads); i++ {
		start := i * BytesPrThread
		end := start + BytesPrThread - 1
		if i == int64(a.threads)-1 {
			end = ContentSize - 1 // Asegurar que el último worker cubra hasta el final
		}
		name := a.Filename + "." + ffmpeg.Int64toString(i)
		names = append(names, name)
		go func() {
			err := ParallelWorker(a, name, a.Url, start, end)
			if err != nil {
				dlerrors <- err
			}
		}()
	}
	for {
		time.Sleep(100 * time.Millisecond)
		dlmutex.Lock()
		if a.running == 0 {
			defer dlmutex.Unlock()
			break
		}
		dlmutex.Unlock()
	}
	return JoinFiles(a.Filename, names)
}
