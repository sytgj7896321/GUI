package main

import (
	"GUI/X11/pics"
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
	"time"
)

var autoFlag bool

func main() {
	go pics.MakeCache()
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

	//Settings
	tFloat := 5.0
	tData := binding.BindFloat(&tFloat)
	tLabel := widget.NewLabelWithData(binding.FloatToStringWithFormat(tData, "Refresh Interval: %0.0fs"))
	tSlide := widget.NewSliderWithData(5, 120, tData)
	tSlide.SetValue(30)

	autoRefresh := widget.NewCheck("Auto Refresh", func(value bool) {
		if value {
			log.Println("Auto Refresh On")
			autoFlag = true
			tSlide.Hide()
			go refreshTick(tData)
		} else {
			log.Println("Auto Refresh Off")
			autoFlag = false
			tSlide.Show()
		}
	})

	f := 0.1
	data := binding.BindFloat(&f)
	//bar := widget.NewProgressBarWithData(data)

	downloadList := binding.BindFloatList(&[]float64{0.6, 0.7})
	list := widget.NewListWithData(downloadList,
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil, nil, widget.NewButton("+", nil),
				widget.NewLabel("item x.y"))
		},
		func(item binding.DataItem, obj fyne.CanvasObject) {
			f := item.(binding.Float)
			text := obj.(*fyne.Container).Objects[0].(*widget.Label)
			text.Bind(binding.FloatToStringWithFormat(f, "pic %0.1f"))

			btn := obj.(*fyne.Container).Objects[1].(*widget.Button)
			btn.OnTapped = func() {
				val, _ := f.Get()
				_ = f.Set(val + 1)
			}
		})

	autoSave := widget.NewCheck("Auto Save Original Pictures to Local Directory After Refresh", func(value bool) {
		if value {
			f = f + 0.1
			data.Reload()
		} else {
			f = f + 0.1
			data.Reload()
		}

	})
	autoSave.Enable()

	currentPath := widget.NewLabel("Local Save Directory: ")
	homeDir, _ := os.UserHomeDir()
	err := createPath(homeDir + "/Pics")
	if err != nil {
		log.Println("Can not create directory in Home Directory, please choose a directory by yourself")
	} else {
		currentPath.Text = "Local Save Directory: " + homeDir + "/Pics"
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
			currentPath.Text = "Local Save Directory: " + strings.TrimPrefix(list.String(), "file://")
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
				closeBtn),
		),
		container.NewTabItemWithIcon(
			"Tasks",
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
		log.Println("Lifecycle: Started")
	})
	fyne.CurrentApp().Lifecycle().SetOnStopped(func() {
		log.Println("Lifecycle: Stopped")
	})
	fyne.CurrentApp().Lifecycle().SetOnEnteredForeground(func() {
		log.Println("Lifecycle: Entered Foreground")
	})
	fyne.CurrentApp().Lifecycle().SetOnExitedForeground(func() {
		log.Println("Lifecycle: Exited Foreground")
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
		if autoFlag {
			pics.RefreshAll()
		} else {
			break
		}
	}
}
