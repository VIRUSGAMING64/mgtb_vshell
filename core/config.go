package core

type Proc struct {
	maxprocs     int
	maxThreads   int
	maxDownloads int
}

var Config *Proc = new(Proc)
