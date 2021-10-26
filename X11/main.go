package main

import (
	"GUI/X11/pics"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/cmd/fyne_settings/settings"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"log"
	"net/url"
)

func main() {
	go pics.MakeCache()
	myApp := app.NewWithID("NewApp")
	logLifecycle()
	//img := <-pics.imageChan
	//image := canvas.NewImageFromReader(img.ImagData, img.Id)
	//myApp.SetIcon(image.Resource)
	mainWindow := myApp.NewWindow("Wallpaper Tool")
	mainWindow.SetMaster()
	mainWindow.Resize(fyne.NewSize(500, 400))

	mainWindow.SetMainMenu(makeMenu(myApp, mainWindow))

	captureBtn := widget.NewButton("Open Capture Window", func() {
		func() {
			pics.CapturePic()
			myApp.SendNotification(&fyne.Notification{
				Title:   "Wallpaper Tool",
				Content: "New Capture Window Opened",
			})
		}()
	})

	refreshBtn := widget.NewButton("Refresh", func() {
		pics.RefreshAll()
	})

	closeBtn := widget.NewButton("Close All Pictures", func() {
		pics.CloseAllWindows()
	})
	bugURL, _ := url.Parse("https://github.com/sytgj7896321/GUI/issues/new")

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon(
			"Home",
			theme.HomeIcon(),
			container.NewHScroll(container.NewVBox(captureBtn, refreshBtn, closeBtn)),
		),
		container.NewTabItemWithIcon(
			"Favourite",
			theme.ListIcon(),
			container.NewHScroll(container.NewVBox(widget.NewLabel("TODO"))),
		),
		container.NewTabItemWithIcon(
			"Settings",
			theme.SettingsIcon(),
			container.NewHScroll(container.NewVBox(widget.NewLabel("TODO"))),
		),
		container.NewTabItemWithIcon(
			"Help",
			theme.HelpIcon(),
			container.NewHScroll(container.NewVBox(widget.NewHyperlink("Report a bug", bugURL))),
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

func makeMenu(app fyne.App, win fyne.Window) *fyne.MainMenu {
	newItem := fyne.NewMenuItem("New", nil)
	checkedItem := fyne.NewMenuItem("Auto Capture", nil)
	checkedItem.Checked = true
	checkedItem.Action = func() {
		checkedItem.Checked = !checkedItem.Checked
		if checkedItem.Checked {
			log.Println("Auto Capture On")
		} else {
			log.Println("Auto Capture Off")
		}
	}
	newItem.ChildMenu = fyne.NewMenu("",
		fyne.NewMenuItem("Window", func() {
			go func() {
				pics.CapturePic()
				app.SendNotification(&fyne.Notification{
					Title:   "Wallpaper Tool",
					Content: "Open Capture Window",
				})
			}()
		}),
	)
	settingsItem := fyne.NewMenuItem("Settings", func() {
		w := app.NewWindow("Fyne Settings")
		w.SetContent(settings.NewSettings().LoadAppearanceScreen(w))
		w.Resize(fyne.NewSize(480, 480))
		w.Show()
	})

	cutItem := fyne.NewMenuItem("Cut", func() {
		shortcutFocused(&fyne.ShortcutCut{
			Clipboard: win.Clipboard(),
		}, win)
	})
	copyItem := fyne.NewMenuItem("Copy", func() {
		shortcutFocused(&fyne.ShortcutCopy{
			Clipboard: win.Clipboard(),
		}, win)
	})
	pasteItem := fyne.NewMenuItem("Paste", func() {
		shortcutFocused(&fyne.ShortcutPaste{
			Clipboard: win.Clipboard(),
		}, win)
	})
	findItem := fyne.NewMenuItem("Find", func() { fmt.Println("Menu Find") })

	helpMenu := fyne.NewMenu("Help",
		fyne.NewMenuItem("Documentation", func() {
			u, _ := url.Parse("https://developer.fyne.io")
			_ = app.OpenURL(u)
		}),
		fyne.NewMenuItem("Support", func() {
			u, _ := url.Parse("https://fyne.io/support/")
			_ = app.OpenURL(u)
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Sponsor", func() {
			u, _ := url.Parse("https://fyne.io/sponsor/")
			_ = app.OpenURL(u)
		}))

	// app quit item will be appended to our first (File) menu
	file := fyne.NewMenu("Home", newItem, checkedItem)
	if !fyne.CurrentDevice().IsMobile() {
		file.Items = append(file.Items, fyne.NewMenuItemSeparator(), settingsItem)
	}
	return fyne.NewMainMenu(
		file,
		fyne.NewMenu("Edit", cutItem, copyItem, pasteItem, fyne.NewMenuItemSeparator(), findItem),
		helpMenu,
	)
}

func shortcutFocused(s fyne.Shortcut, w fyne.Window) {
	if focused, ok := w.Canvas().Focused().(fyne.Shortcutable); ok {
		focused.TypedShortcut(s)
	}
}
