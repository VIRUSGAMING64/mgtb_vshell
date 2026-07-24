package ffmpeg

func Removereps(arr []string) []string {
	var news []string
	for i := 0; i < len(arr); i++ {
		if IsIn(arr[i], news) == -1 {
			news = append(news, arr[i])
		}
	}
	return news
}

func String2Int64(data string) int64 {
	d := int64(0)
	c := int64(1)
	for i := len(data) - 1; i >= 0; i-- {
		if -1 == StringIsIn(data[i], "0123456789") {
			return -1
		}
		d = d + int64(data[i]-'0')*c
		c = c * 10
	}
	return d
}

func Int64toString(c int64) string {
	if c == 0 {
		return "0"
	}
	a := string("")
	for x := c; x != 0; x = x / 10 {
		a = string(byte(x%10)+'0') + a
	}
	return a
}

func IsIn(elem string, arr []string) int {
	for i := 0; i < len(arr); i++ {
		if elem == arr[i] {
			return i
		}
	}
	return -1
}

func StringIsIn(elem byte, arr string) int {
	for i := 0; i < len(arr); i++ {
		if elem == arr[i] {
			return i
		}
	}
	return -1
}
