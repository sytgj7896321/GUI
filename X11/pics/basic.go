package pics

import (
	"GUI/X11/fetcher"
	"bytes"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
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
	Win      fyne.Window
	Refresh  func()
	Position int
	FullLink string
}

const (
	random   = "https://wallhaven.cc/random"
	selector = "#thumbs > section > ul"
	small    = "https://th.wallhaven.cc/small/"
	full     = "https://w.wallhaven.cc/full/"
)

var (
	imageChan      = make(chan *Image, 96)
	CaptureWindows []*Window
	AutoSaveFlag   = false
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
	win := fyne.CurrentApp().NewWindow("Picture")
	win.Resize(fyne.NewSize(300, 200))
	img := <-imageChan
	image := canvas.NewImageFromReader(img.ImagData, img.Id)
	image.Resize(fyne.NewSize(300, 200))
	win.SetContent(image)
	win.Show()
	myWin := new(Window)
	myWin.Win = win
	myWin.FullLink = img.Full
	myWin.Refresh = func() {
		if AutoSaveFlag {

		}
		img = <-imageChan
		myWin.FullLink = img.Full
		image = canvas.NewImageFromReader(img.ImagData, img.Id)
		image.Resize(fyne.NewSize(300, 200))
		myWin.Win.Canvas().SetContent(image)
	}
	CaptureWindows = append(CaptureWindows, myWin)
	myWin.Position = len(CaptureWindows) - 1
	win.SetCloseIntercept(func() {
		win.Close()
		for _, w := range CaptureWindows[myWin.Position+1:] {
			w.Position--
		}
		CaptureWindows = append(CaptureWindows[:myWin.Position], CaptureWindows[myWin.Position+1:]...)
	})
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
				go downloadSmallImage(idList[i])
			}
		}
	}
}

func GetLength() string {
	return strconv.Itoa(len(CaptureWindows))
}

func downloadSmallImage(id string) {
	img := new(Image)
	img.Id = id
	img.Small = small + string([]byte(id)[:2]) + "/" + id + ".jpg"
	img.Full = full + string([]byte(id)[:2]) + "/wallhaven-" + id
	body, err := fetcher.Fetch(img.Small)
	if err != nil {
		log.Println(err)
	}
	img.ImagData = bytes.NewReader(body)
	imageChan <- img
}
