package pics

import (
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
	ImgType  string
	Original string
	ImagData io.Reader
}

type Window struct {
	Win      fyne.Window
	Refresh  func()
	Position int
}

const (
	//random   = "https://wallhaven.cc/random"
	random = "https://wallhaven.cc/search?categories=010&purity=010&sorting=random"
	//random   = "https://wallhaven.cc/search?categories=111&purity=110&sorting=random"
	selector = "#thumbs > section > ul"
	small    = "https://th.wallhaven.cc/small/"
	Full     = "https://w.wallhaven.cc/full/"
)

var (
	imageChan      = make(chan *Image, 96)
	CaptureWindows []*Window
	AutoSaveFlag   = false
)

func CloseAllWindows() {
	for _, w := range CaptureWindows {
		c := w.Win
		go c.Close()
	}
	CaptureWindows = CaptureWindows[len(CaptureWindows):]
}

func RefreshAll() {
	for _, r := range CaptureWindows {
		go r.Refresh()
	}
}

func CapturePic() {
	img := <-imageChan
	win := fyne.CurrentApp().NewWindow("Picture")
	win.Resize(fyne.NewSize(300, 200))
	image := canvas.NewImageFromReader(img.ImagData, img.Id)
	image.Resize(fyne.NewSize(300, 200))
	win.SetContent(image)
	win.Show()
	myWin := new(Window)
	myWin.Win = win
	myWin.Refresh = func() {
		if AutoSaveFlag {
			go func() {
				task := new(Task)
				task.ImageId = &img.Id
				task.Link = &img.Original
				TaskList = append(TaskList, task)
			}()
		}
		img = <-imageChan
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
			body, err := Fetch(random)
			if err != nil {
				panic(err)
			}
			dom, _ := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
			idList := dom.Find(selector).Contents().Map(func(i int, selection *goquery.Selection) string {
				e, _ := selection.Children().Attr("data-wallpaper-id")
				return e
			})
			imgTypeList := dom.Find(selector).Contents().Map(func(i int, selection *goquery.Selection) string {
				e := selection.Children().Children().ChildrenFiltered("span.png").Text()
				return e
			})
			for i := range idList {
				go downloadSmallImage(idList[i], imgTypeList[i])
			}
		}
	}
}

func GetLength() string {
	return strconv.Itoa(len(CaptureWindows))
}

func downloadSmallImage(id, imgType string) {
	img := new(Image)
	img.Id = id
	if imgType == "PNG" {
		img.ImgType = ".png"
	} else {
		img.ImgType = ".jpg"
	}
	img.Original = Full + string([]byte(img.Id)[:2]) + "/wallhaven-" + img.Id + img.ImgType
	body, err := Fetch(small + string([]byte(id)[:2]) + "/" + id + ".jpg")
	img.ImagData = bytes.NewReader(body)
	if err != nil {
		log.Println(err)
		img.ImagData = nil
	}
	imageChan <- img
}
