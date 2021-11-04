package main

import (
	"GUI/X11/pics"
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
	"strings"
	"time"
)

var (
	autoFlag     = make(chan bool, 1)
	downloadList = binding.BindFloatList(&[]float64{})
	list         = widget.NewListWithData(
		downloadList,
		nil,
		func(item binding.DataItem, obj fyne.CanvasObject) {
			f := item.(binding.Float)
			bar := obj.(*fyne.Container).Objects[0].(*widget.ProgressBar)
			bar.Bind(f)
		})
)

func main() {
	go pics.MakeCache()
	pics.DownloadOriginal()
	myApp := app.NewWithID("NewApp")
	icon, _ := fyne.LoadResourceFromPath("/usr/local/share/pixmaps/WallpaperTool.png")
	myApp.SetIcon(icon)
	lifecycle()
	mainWindow := myApp.NewWindow("Wallpaper Tool")
	mainWindow.SetMaster()
	mainWindow.Resize(fyne.NewSize(600, 350))

	//Home
	captureBtn := widget.NewButton("Open New Capture Window", func() {
		if len(pics.CaptureWindows) < 24 {
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
		Content: "Total Windows Opened: " + pics.GetLength(),
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
