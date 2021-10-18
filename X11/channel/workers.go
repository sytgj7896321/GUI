package channel

import (
	"fmt"
	"os"
	"sync"
)

type Windows struct {
	open  func()
	close []func()
}

func CreateWindows(wg *sync.WaitGroup, rp, fp *os.File) Windows {
	w := Windows{
		open:  nil,
		close: nil,
	}
	go CloseWindows(w, rp, fp)
	return w
}

func CloseWindows(w Windows, rp, fp *os.File) {
	for i := range w.close {
		fmt.Println(i)
	}
}
