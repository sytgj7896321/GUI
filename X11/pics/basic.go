package pics

import (
	"GUI/X11/fetcher"
	"bytes"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
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

type Window struct {
	Win     fyne.Window
	Refresh func()
}

const (
	random   = "https://wallhaven.cc/random"
	selector = "#thumbs > section > ul"
	small    = "https://th.wallhaven.cc/small/"
	full     = "https://w.wallhaven.cc/full/"
)

var (
	imageChan      = make(chan Image, 96)
	captureWindows []*Window
)

func CloseAllWindows() {
	for _, w := range captureWindows {
		c := w.Win
		c.Close()
	}
	captureWindows = captureWindows[len(captureWindows):]
}

func CapturePic() {
	win := fyne.CurrentApp().NewWindow("Picture")
	win.Resize(fyne.NewSize(450, 300))
	img := <-imageChan
	image := canvas.NewImageFromReader(img.ImagData, img.Id)
	image.Resize(fyne.NewSize(450, 300))
	//toolbar := widget.NewToolbar(
	//	widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
	//		img = <-imageChan
	//		image = canvas.NewImageFromReader(img.ImagData, img.Id)
	//		content := container.NewBorder(toolBar(win), nil, nil, nil, image)
	//		win.SetContent(content)
	//	}),
	//	widget.NewToolbarAction(theme.ViewFullScreenIcon(), func() {
	//		img = <-imageChan
	//		image = canvas.NewImageFromReader(img.ImagData, img.Id)
	//		win.SetContent(image)
	//		win.SetFullScreen(true)
	//	}),
	//	widget.NewToolbarAction(theme.ContentCopyIcon(), func() {
	//		fmt.Println("Copy to Clipboard")
	//	}),
	//	widget.NewToolbarAction(theme.DownloadIcon(), func() {
	//		fmt.Println("Download to Local")
	//	}),
	//)
	//content := container.NewBorder(toolbar, nil, nil, nil, image)
	win.SetContent(image)
	win.Show()
	myWin := new(Window)
	myWin.Win = win
	myWin.Refresh = func() {
		img = <-imageChan
		image = canvas.NewImageFromReader(img.ImagData, img.Id)
		image.Resize(fyne.NewSize(450, 300))
		myWin.Win.Canvas().SetContent(image)
	}
	captureWindows = append(captureWindows, myWin)

}

func RefreshAll() {
	for _, r := range captureWindows {
		r.Refresh()
	}
}

func MakeCache() {
	for true {
		if len(imageChan) <= 48 {
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
				go downloadImage(idList[i])
			}
		}
	}
}

func downloadImage(id string) {
	img := Image{
		Id:    id,
		Small: small + string([]byte(id)[:2]) + "/" + id + ".jpg",
		Full:  full + string([]byte(id)[:2]) + "/wallhaven-" + id + ".jpg",
	}
	//TODO ProgressBar
	body, err := fetcher.Fetch(img.Small)
	if err != nil {
		return
	}
	img.ImagData = bytes.NewReader(body)
	imageChan <- img
}

func toolBar(myWindow fyne.Window) fyne.CanvasObject {
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
			img := <-imageChan
			image := canvas.NewImageFromReader(img.ImagData, img.Id)
			content := container.NewBorder(toolBar(myWindow), nil, nil, nil, image)
			myWindow.SetContent(content)
		}),
		widget.NewToolbarAction(theme.ViewFullScreenIcon(), func() {
			img := <-imageChan
			image := canvas.NewImageFromReader(img.ImagData, img.Id)
			myWindow.SetContent(image)
			myWindow.SetFullScreen(true)
		}),
		widget.NewToolbarAction(theme.ContentCopyIcon(), func() {
			fmt.Println("Copy to Clipboard")
		}),
		widget.NewToolbarAction(theme.DownloadIcon(), func() {
			fmt.Println("Download to Local")
		}),
	)
	return toolbar
}
