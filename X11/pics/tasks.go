package pics

import (
	"github.com/cavaliercoder/grab"
)

type Task struct {
	ImageId *string
	Link    *string
	Percent float64
	Done    func()
	Cancel  func()
	Retry   func()
}

var (
	Routines = 8
	In       = make(chan *grab.Request, Routines)
	Out      = make(chan *grab.Response, Routines)
)

func DownloadOriginal() {
	client := grab.NewClient()
	for i := 0; i < Routines; i++ {
		go client.DoChannel(In, Out)
	}
}
