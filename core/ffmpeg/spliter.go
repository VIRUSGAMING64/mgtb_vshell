package ffmpeg

import (
	"fmt"
	"os/exec"
)

/*
ffmpeg -i "test.mp4" -c copy -map 0 -f segment -segment_time 120 "test.mp4.%03d.mp4"
*/
func SplitInto(vid string, duration int64) error {
	ff := exec.Command(
		"ffmpeg",
		"-i", vid,
		"-c", "copy",
		"-map", "0",
		"-f", "segment",
		"-segment_time", Int64toString(duration),
		""+vid+"."+"%"+"03d.mp4",
	)
	fmt.Println(ff.Args)
	out, err := ff.CombinedOutput()
	fmt.Println(string(out))
	return err
}
