package pics

import (
	"GUI/X11/fetcher"
	"bytes"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
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

var (
	imageChan = make(chan Image, 24)
)

func CloseAllWindows(windows []fyne.Window) {
	for _, w := range windows {
		w.Close()
	}
}

func CapturePic(app fyne.App) fyne.Window {
	img := <-imageChan
	myWindow := app.NewWindow(img.Id)
	myWindow.Resize(fyne.NewSize(450, 350))
	myCanvas := myWindow.Canvas()
	image := canvas.NewImageFromReader(img.ImagData, img.Id)
	myCanvas.SetContent(image)

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
			log.Println("New document")
		}),
		widget.NewToolbarAction(theme.ContentCutIcon(), func() {}),
		widget.NewToolbarAction(theme.ContentCopyIcon(), func() {}),
		widget.NewToolbarAction(theme.ContentPasteIcon(), func() {}),
		widget.NewToolbarAction(theme.HelpIcon(), func() {
			log.Println("Display help")
		}),
	)
	content := container.NewBorder(toolbar, nil, nil, nil, image)
	myWindow.SetContent(content)
	myWindow.Show()

	return myWindow
}

func MakeCache() {
	for true {
		if len(imageChan) <= 8 {
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
				go createImage(idList[i])
			}
		}
	}
}

func createImage(id string) {
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
	imageChan <- img
}
