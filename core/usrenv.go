package core

import (
	"fmt"
	"govshell/core/ffmpeg"
	"os"
	"path"
	"strings"
)

var ROOT = "roots/"

type User struct {
	Root   string `json:"root"`
	Path   string `json:"path"`
	Name   string `json:"name"`
	Id     int64  `json:"id"`
	Status int    `json:"status"`
}

var Users map[int64]*User = make(map[int64]*User)

func NewUser(id int64) *User {
	os.MkdirAll(ROOT, 0644)
	os.MkdirAll(path.Join(ROOT, ffmpeg.Int64toString(int64(id))), 0644)
	usr := User{Id: id}
	usr.Path = path.Join(ROOT, ffmpeg.Int64toString(int64(id)))
	usr.Root = usr.Path
	return &usr
}

func _back(a *User) error {
	if a.Path == a.Root {
		return fmt.Errorf("current dir is root")
	}
	pth, e := strings.CutSuffix(a.Path, path.Base(a.Path))
	if e {
		a.Path = pth
	}
	return nil
}

func (a *User) Mkdir(dir string) error {
	ndir := path.Join(a.Path, dir)
	return os.MkdirAll(ndir, 0644)
}

func (a *User) Chdir(dir string) error {
	if dir == "." {
		return nil
	}
	if dir == ".." {
		return _back(a)
	}
	pth := path.Join(a.Path, dir)
	fd, e := os.Open(pth)
	if e != nil {
		return e
	}
	defer fd.Close()
	info, e := fd.Stat()
	if e != nil {
		return e
	}
	if info.IsDir() == false {
		return fmt.Errorf("Not directory")
	}
	a.Path = pth
	return nil
}

func (a *User) ListDir() ([]os.DirEntry, error) {
	fl, e := os.Open(a.Path)
	if e != nil {
		return nil, e
	}
	defer fl.Close()
	dirs, e := os.ReadDir(a.Path)
	if e != nil {
		return nil, e
	}
	return dirs, nil
}

func (a *User) Rm(name string) error {
	ruta := path.Join(a.Path, name)
	return os.RemoveAll(ruta)

}
