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
	"log"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	autoFlag = make(chan bool, 1)
)

func main() {
	go pics.MakeCache()
	//go func() {
	//	for t := range pics.TaskList {
	//		pics.AdvancedDownloadOriginal(t)
	//	}
	//}()
	pics.DownloadOriginal()
	myApp := app.NewWithID("NewApp")
	logLifecycle()
	mainWindow := myApp.NewWindow("Wallpaper Tool")
	mainWindow.SetMaster()
	mainWindow.Resize(fyne.NewSize(600, 400))

	//Home
	captureBtn := widget.NewButton("New Capture Window", func() {
		if len(pics.CaptureWindows) < 24 {
			pics.CapturePic()
			myApp.SendNotification(&fyne.Notification{
				Title:   "Wallpaper Tool",
				Content: "New Capture Window Opened",
			})
			getWindowsNum()
		} else {
			warn := dialog.NewConfirm("Warning", "Too much windows opened\nAre you still want to add another one?(Not recommended)", func(b bool) {
				if b {
					pics.CapturePic()
					myApp.SendNotification(&fyne.Notification{
						Title:   "Wallpaper Tool",
						Content: "New Capture Window Opened",
					})
					getWindowsNum()
				} else {
					log.Println("Cancel Open New Window")
					getWindowsNum()
				}
			}, mainWindow)
			warn.SetDismissText("NO")
			warn.SetConfirmText("YES")
			warn.Show()
		}
	})

	refreshBtn := widget.NewButton("Refresh", func() {
		pics.RefreshAll()
		myApp.SendNotification(&fyne.Notification{
			Title:   "Wallpaper Tool",
			Content: "All Pictures Refreshed",
		})
		getWindowsNum()
	})

	closeBtn := widget.NewButton("Close All Pictures", func() {
		pics.CloseAllWindows()
		myApp.SendNotification(&fyne.Notification{
			Title:   "Wallpaper Tool",
			Content: "All Windows Closed",
		})
		getWindowsNum()
	})

	//Tasks
	downloadList := binding.BindFloatList(&[]float64{})
	tmpBtn := widget.NewButton("tmp", func() {
		downloadList.Append(float64(downloadList.Length()+1) / 10)

	})
	list := widget.NewListWithData(
		downloadList,
		func() fyne.CanvasObject {
			label := widget.NewLabel("unknown")
			//if len(pics.TaskList) > 0 {
			//	label.SetText(pics.TaskList[0].ImageId)
			//	pics.TaskList = pics.TaskList[1:]
			//}
			bar := widget.NewProgressBar()
			menu := widget.NewMenu(
				fyne.NewMenu("",
					//TODO
					fyne.NewMenuItem("remove", func() {
						fmt.Println("TODO")
					}),
					fyne.NewMenuItem("retry", func() {
						fmt.Println("TODO")
					}),
				))
			return container.NewBorder(nil, nil, label, menu, bar)
		},
		func(item binding.DataItem, obj fyne.CanvasObject) {
			f := item.(binding.Float)
			bar := obj.(*fyne.Container).Objects[0].(*widget.ProgressBar)
			bar.Bind(f)
		})

	//Settings
	tFloat := 5.0
	tData := binding.BindFloat(&tFloat)
	tLabel := widget.NewLabelWithData(binding.FloatToStringWithFormat(tData, "Refresh Interval: %0.0fs"))
	tSlide := widget.NewSliderWithData(5, 60, tData)
	tSlide.SetValue(30)

	autoSave := widget.NewCheck("Auto Save Original Pictures to Local Directory After Refresh", func(value bool) {
		if value {
			pics.AutoSaveFlag = true
		} else {
			pics.AutoSaveFlag = false
		}
	})

	autoRefresh := widget.NewCheck("Auto Refresh", func(value bool) {
		var wg sync.WaitGroup
		if value {
			log.Println("Auto Refresh On")
			tSlide.Hide()
			wg.Add(1)
			go refreshTick(tData)
		} else {
			log.Println("Auto Refresh Off")
			autoFlag <- false
			tSlide.Show()
		}
	})

	currentPath := widget.NewLabel("Local Save Directory: ")
	pics.LocalSaveDirectory, _ = os.UserHomeDir()
	pics.LocalSaveDirectory = pics.LocalSaveDirectory + "/Pics"
	err := createPath(pics.LocalSaveDirectory)
	if err != nil {
		log.Println("Can not create directory in Home Directory, please choose a directory by yourself")
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
				log.Println("Cancelled")
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
				refreshBtn,
				closeBtn,
				tmpBtn),
		),
		container.NewTabItemWithIcon(
			"Download",
			theme.DownloadIcon(),
			container.NewGridWithColumns(
				1,
				list),
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

func logLifecycle() {
	fyne.CurrentApp().Lifecycle().SetOnStarted(func() {
		log.Println("Wallpaper Tool: Started")
	})
	fyne.CurrentApp().Lifecycle().SetOnStopped(func() {
		log.Println("Wallpaper Tool: Stopped")
	})
	fyne.CurrentApp().Lifecycle().SetOnEnteredForeground(func() {
		log.Println("Wallpaper Tool: Entered Foreground")
	})
	fyne.CurrentApp().Lifecycle().SetOnExitedForeground(func() {
		log.Println("Wallpaper Tool: Exited Foreground")
	})
}

func getWindowsNum() {
	log.Println("Total Windows Opened: " + pics.GetLength())
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
		log.Println("Directory '" + path + "' created")
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
