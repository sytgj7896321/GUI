package pics

import (
	"GUI/X11/fetcher"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/storage/repository"
	"github.com/PuerkitoBio/goquery"
	"io"
	"strings"
)

type Image struct {
	Small    string
	Full     string
	ImagData *io.ReadCloser
}

const (
	random   = "https://wallhaven.cc/random"
	selector = "#thumbs > section > ul"
	small    = "https://th.wallhaven.cc/small/"
	full     = "https://w.wallhaven.cc/full/"
)

var imageChan = make(chan Image, 24)

func CloseAllWindows(windows []fyne.Window) {
	for _, w := range windows {
		w.Close()
	}
}

func CapturePic(app fyne.App) fyne.Window {
	myWindow := app.NewWindow("Picture")
	myCanvas := myWindow.Canvas()
	//TODO
	uri, _ := repository.ParseURI("PicUri")
	myCanvas.SetContent(canvas.NewImageFromURI(uri))
	myWindow.Show()
	return myWindow
}

func MakeCache() {
	body, _ := fetcher.Fetch(random)
	dom, _ := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	idList := dom.Find(selector).Contents().Map(func(i int, selection *goquery.Selection) string {
		e, _ := selection.Children().Attr("data-wallpaper-id")
		return e
	})
	for i := range idList {
		imageChan <- createImage(idList[i])
	}
	for true {
		fmt.Printf("%+v\n", <-imageChan)
	}
}

func createImage(id string) Image {
	return Image{
		Small:    small + string([]byte(id)[:2]) + "/" + id + ".jpg",
		Full:     full + string([]byte(id)[:2]) + "/wallhaven-" + id + ".jpg",
		ImagData: nil,
	}
}
