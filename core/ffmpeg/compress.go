package ffmpeg

import (
	"os"
	"os/exec"
	"time"
)

/*
Requires name of video to compress
*/
type Callback func(part, total int64)

func Compress(vid string, power bool, fn Callback) error {
	out := vid + ".cmp.mp4"
	th := int64(1)
	arg := "pools=1:frame-threads=1:pmode=0:pme=0:rc-lookahead=10"
	if power {
		th = 8
		arg = ""
	}
	threads := Int64toString(th)
	ff := exec.Command(
		"ffmpeg",
		"-i", vid,
		"-threads", threads,
		"-c:v", "libx265",
		"-preset", "medium",
		//"-crf", "18",
		"-x265-params",
		arg,
		"-c:a", "copy",
		out,
	)
	var err error = nil
	go func() {
		for ff.ProcessState.ExitCode() == -1 {
			time.Sleep(3 * time.Second)
			f1, e1 := os.OpenFile(vid, os.O_RDONLY, 0644)
			f2, e2 := os.OpenFile(out, os.O_RDONLY, 0644)
			defer f1.Close()
			defer f2.Close()
			if e1 != nil || e2 != nil || err != nil {
				return
			}
			s1, _ := f1.Stat()
			s2, _ := f2.Stat()
			if fn != nil {
				fn(s2.Size(), s1.Size())
			}
		}
	}()
	err = ff.Run()
	return err
}
