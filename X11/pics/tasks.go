package pics

import (
	"github.com/cavaliercoder/grab"
)

var (
	Routines = 4
	In       = make(chan *grab.Request, 24)
	Out      = make(chan *grab.Response, Routines)
)

func DownloadOriginal() {
	client := grab.NewClient()
	for i := 0; i < Routines; i++ {
		go client.DoChannel(In, Out)
	}
}
