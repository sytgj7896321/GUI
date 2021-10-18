package pics

import (
	"GUI/X11/fetcher"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/storage/repository"
)

var (
	root     = "https://wallhave.cc/random"
	selector = "#thumbs > section > ul"
	uriList  []string
)

func CloseAllWindows(windows []fyne.Window) {
	for _, w := range windows {
		w.Close()
	}
}

func CapturePic(app fyne.App) fyne.Window {
	myWindow := app.NewWindow("Picture")
	myCanvas := myWindow.Canvas()

	uri, _ := repository.ParseURI("https://w.wallhaven.cc/full/8o/wallhaven-8o96xo.jpg")
	myCanvas.SetContent(canvas.NewImageFromURI(uri))
	myWindow.Show()
	return myWindow
}

func MakeCache() {
	fmt.Println(fetcher.Fetch(root))
}
