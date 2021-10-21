package pics

import (
	"GUI/X11/fetcher"
	"bytes"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"github.com/PuerkitoBio/goquery"
	"io"
	"strings"
)

type Image struct {
	Id       string
	Small    string
	Full     string
	ImagData io.Reader
}

const (
	random   = "https://wallhaven.cc/random"
	selector = "#thumbs > section > ul"
	small    = "https://th.wallhaven.cc/small/"
	full     = "https://w.wallhaven.cc/full/"
)

var imageChan = make(chan Image, 48)

func CloseAllWindows(windows []fyne.Window) {
	for _, w := range windows {
		w.Close()
	}
}

func CapturePic(app fyne.App) fyne.Window {
	img := <-imageChan
	myWindow := app.NewWindow(img.Id)
	myWindow.Resize(fyne.Size{
		Width:  300,
		Height: 200,
	})
	myCanvas := myWindow.Canvas()
	//TODO
	myCanvas.SetContent(canvas.NewImageFromReader(img.ImagData, img.Id))
	myWindow.Show()
	return myWindow
}

func MakeCache() {
	for true {
		body, err := fetcher.Fetch(random)
		if err != nil {
			panic(err)
		}
		dom, _ := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
		idList := dom.Find(selector).Contents().Map(func(i int, selection *goquery.Selection) string {
			e, _ := selection.Children().Attr("data-wallpaper-id")
			return e
		})
		for i := range idList {
			imageChan <- createImage(idList[i])
		}
	}
}

func createImage(id string) Image {
	img := Image{
		Id:    id,
		Small: small + string([]byte(id)[:2]) + "/" + id + ".jpg",
		Full:  full + string([]byte(id)[:2]) + "/wallhaven-" + id + ".jpg",
	}
	data, err := fetcher.Fetch(img.Small)
	if err != nil {
		panic(err)
	}
	reader := bytes.NewReader(data)
	img.ImagData = reader
	return img
}
