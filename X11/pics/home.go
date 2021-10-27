package pics

import (
	"GUI/X11/fetcher"
	"bytes"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"github.com/PuerkitoBio/goquery"
	"io"
	"strconv"
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
	CaptureWindows []*Window
)

func CloseAllWindows() {
	for _, w := range CaptureWindows {
		c := w.Win
		c.Close()
	}
	CaptureWindows = CaptureWindows[len(CaptureWindows):]
}

func RefreshAll() {
	for _, r := range CaptureWindows {
		go r.Refresh()
	}
}

func CapturePic() {
	win := fyne.CurrentApp().NewWindow("Picture" + GetLength(true))
	win.Resize(fyne.NewSize(300, 200))
	img := <-imageChan
	image := canvas.NewImageFromReader(img.ImagData, img.Id)
	image.Resize(fyne.NewSize(300, 200))
	win.SetContent(image)
	win.Show()
	myWin := new(Window)
	myWin.Win = win
	myWin.Refresh = func() {
		img = <-imageChan
		image = canvas.NewImageFromReader(img.ImagData, img.Id)
		image.Resize(fyne.NewSize(300, 200))
		myWin.Win.Canvas().SetContent(image)
	}
	CaptureWindows = append(CaptureWindows, myWin)

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

func GetLength(b bool) string {
	if b {
		return strconv.Itoa(len(CaptureWindows) + 1)
	}
	return strconv.Itoa(len(CaptureWindows))
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
