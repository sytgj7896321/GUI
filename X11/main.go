package main

import (
	"GUI/X11/pics"
	"encoding/base64"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/cavaliercoder/grab"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	autoFlag     = make(chan bool, 1)
	downloadList = binding.BindFloatList(&[]float64{})
	iconPng      = `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
D+atpJXmvwGH/alhp+KhwAAAAABJRU5ErkJggg==`
)

type Icon struct{}

func (i Icon) Name() string {
	return "Icon"
}

func (i Icon) Content() []byte {
	decodeString, _ := base64.StdEncoding.DecodeString(iconPng)
	return decodeString
}

func main() {
	go pics.MakeCache()
	pics.DownloadOriginal()
	myApp := app.NewWithID("WallpaperTool")
	myApp.SetIcon(new(Icon))
	lifecycle()
	mainWindow := myApp.NewWindow("Wallpaper Tool")
	mainWindow.SetMaster()
	mainWindow.Resize(fyne.NewSize(600, 350))

	//Home
	captureBtn := widget.NewButton("Open New Capture Window", func() {
		if pics.GetLength() < 24 {
			pics.CapturePic()
		} else {
			warn := dialog.NewConfirm("Warning", "Too much windows opened\nAre you still want to add another one?(Not recommended)", func(b bool) {
				if b {
					pics.CapturePic()
				} else {
					return
				}
			}, mainWindow)
			warn.SetDismissText("NO")
			warn.SetConfirmText("YES")
			warn.Show()
		}
	})

	countBtn := widget.NewButton("Count Windows", func() {
		getWindowsNum()
	})

	refreshBtn := widget.NewButton("Refresh Capture Windows Content", func() {
		pics.RefreshAll()
		myApp.SendNotification(&fyne.Notification{
			Title:   "Wallpaper Tool",
			Content: "All Pictures Refreshed",
		})
	})

	closeBtn := widget.NewButton("Close All Pictures", func() {
		pics.CloseAllWindows()
		myApp.SendNotification(&fyne.Notification{
			Title:   "Wallpaper Tool",
			Content: "All Windows Closed",
		})
	})

	//Tasks
	clearBtn := widget.NewButton("Clear Task List", func() {
		err := downloadList.Set([]float64{})
		if err != nil {
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   "Wallpaper Tool",
				Content: "Clear Task List Failed",
			})
		} else {
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   "Wallpaper Tool",
				Content: "Task List Cleared",
			})
		}
	})
	var list = widget.NewListWithData(
		downloadList,
		nil,
		func(item binding.DataItem, obj fyne.CanvasObject) {
			f := item.(binding.Float)
			bar := obj.(*fyne.Container).Objects[0].(*widget.ProgressBar)
			bar.Bind(f)
		})
	downloadContainer := container.NewVSplit(list, clearBtn)
	downloadContainer.SetOffset(0.8)

	go GetOutData(downloadList, list)

	//Settings
	tFloat := 5.0
	tData := binding.BindFloat(&tFloat)
	tLabel := widget.NewLabelWithData(binding.FloatToStringWithFormat(tData, "Refresh Interval: %0.0fs"))
	tSlide := widget.NewSliderWithData(15, 120, tData)
	tSlide.SetValue(60)

	autoSave := widget.NewCheck("Auto Save Original Pictures to Local Directory After Refresh", func(value bool) {
		if value {
			pics.AutoSaveFlag = true
		} else {
			pics.AutoSaveFlag = false
		}
	})

	autoRefresh := widget.NewCheck("Auto Refresh", func(value bool) {
		if value {
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   "Wallpaper Tool",
				Content: "Auto Refresh On",
			})
			tSlide.Hide()
			go refreshTick(tData)
		} else {
			autoFlag <- false
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   "Wallpaper Tool",
				Content: "Auto Refresh Off",
			})
			tSlide.Show()
		}
	})

	currentPath := widget.NewLabel("Local Save Directory: ")
	pics.LocalSaveDirectory, _ = os.UserHomeDir()
	if runtime.GOOS == "windows" {
		pics.LocalSaveDirectory = pics.LocalSaveDirectory + "\\Pics"
	} else {
		pics.LocalSaveDirectory = pics.LocalSaveDirectory + "/Pics"
	}
	err := createPath(pics.LocalSaveDirectory)
	if err != nil {
		fyne.CurrentApp().SendNotification(&fyne.Notification{
			Title:   "Wallpaper Tool",
			Content: "Can not create directory in Home Directory\nPlease choose a directory by yourself",
		})
	} else {
		currentPath.Text = "Local Save Directory: " + pics.LocalSaveDirectory
		currentPath.Refresh()
	}
	localSavePath := widget.NewButton("Select Local Save Directory", func() {
		dialog.ShowFolderOpen(func(list fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, mainWindow)
				return
			}
			if list == nil {
				return
			}
			pics.LocalSaveDirectory = strings.TrimPrefix(list.String(), "file://")
			currentPath.Text = "Local Save Directory: " + pics.LocalSaveDirectory
			currentPath.Refresh()
		}, mainWindow)
	})

	//Help
	bugURL, _ := url.Parse("https://github.com/sytgj7896321/GUI/issues/new")

	//Create Tabs
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon(
			"Home",
			theme.HomeIcon(),
			container.NewVBox(
				captureBtn,
				countBtn,
				refreshBtn,
				closeBtn),
		),
		container.NewTabItemWithIcon(
			"Download",
			theme.DownloadIcon(),
			downloadContainer,
		),
		container.NewTabItemWithIcon(
			"Settings",
			theme.SettingsIcon(),
			container.NewVBox(
				container.NewGridWithColumns(2, tLabel, tSlide),
				autoRefresh,
				autoSave,
				currentPath,
				localSavePath),
		),
		container.NewTabItemWithIcon(
			"Help",
			theme.HelpIcon(),
			container.NewVBox(
				widget.NewHyperlink("Report a bug", bugURL)),
		),
	)
	tabs.SetTabLocation(container.TabLocationLeading)
	mainWindow.SetContent(tabs)
	mainWindow.Show()
	myApp.Run()
}

func lifecycle() {
	fyne.CurrentApp().Lifecycle().SetOnStarted(func() {
		fyne.CurrentApp().SendNotification(&fyne.Notification{
			Title:   "Wallpaper Tool",
			Content: "Started",
		})
	})
	fyne.CurrentApp().Lifecycle().SetOnStopped(func() {
		fyne.CurrentApp().SendNotification(&fyne.Notification{
			Title:   "Wallpaper Tool",
			Content: "Stopped",
		})
	})
}

func getWindowsNum() {
	fyne.CurrentApp().SendNotification(&fyne.Notification{
		Title:   "Wallpaper Tool",
		Content: "Total Windows Opened: " + strconv.Itoa(pics.GetLength()),
	})
}

func createPath(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		err := os.Mkdir(path, 0755)
		if err != nil {
			return err
		}
		return nil
	}
	return err
}

func refreshTick(t binding.ExternalFloat) {
	for range time.Tick(func(binding.ExternalFloat) time.Duration {
		_ = t.Reload()
		tick, _ := t.Get()
		return time.Duration(tick) * time.Second
	}(t)) {
		select {
		case <-autoFlag:
			return
		default:
			pics.RefreshAll()
		}
	}
}

func GetOutData(downloadList binding.ExternalFloatList, list *widget.List) {
	for {
		select {
		case resp := <-pics.Out:
			list.CreateItem = func() fyne.CanvasObject {
				bar := widget.NewProgressBar()
				bar.TextFormatter = func() string {
					if runtime.GOOS == "windows" {
						return fmt.Sprintf(
							"%s completed %d%%",
							strings.TrimPrefix(resp.Filename, pics.LocalSaveDirectory+"\\"),
							int64(100*bar.Value),
						)
					}
					return fmt.Sprintf(
						"%s completed %d%%",
						strings.TrimPrefix(resp.Filename, pics.LocalSaveDirectory+"/"),
						int64(100*bar.Value),
					)
				}
				return container.NewMax(bar)
			}
			position := operateResponse(resp, downloadList, list)
			tick := time.NewTicker(25 * time.Millisecond)
		Loop:
			for {
				select {
				case <-tick.C:
					_ = downloadList.SetValue(position, resp.Progress())
				case <-resp.Done:
					tick.Stop()
					_ = downloadList.SetValue(position, resp.Progress())
					break Loop
				}
			}
			if err := resp.Err(); err != nil {
				errString := resp.Request.URL().String() + "download failed"
				errWin := fyne.CurrentApp().NewWindow("Error")
				errWin.SetContent(widget.NewTextGridFromString(errString))
				errWin.Resize(fyne.NewSize(float32(len(errString))+10, 50))
				errWin.Show()
				continue
			}
		}
	}
}

func operateResponse(resp *grab.Response, downloadList binding.ExternalFloatList, list *widget.List) int {
	_ = downloadList.Append(resp.Progress())
	list.ScrollToBottom()
	position := downloadList.Length() - 1
	return position
}
