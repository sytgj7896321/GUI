package main

import (
	"GUI/X11/pics"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
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
	mainWindow := myApp.NewWindow("Wallpaper Tool")
	mainWindow.SetMaster()
	mainWindow.Resize(fyne.NewSize(500, 400))

	//mainWindow.SetMainMenu(makeMenu(myApp, mainWindow))

	captureBtn := widget.NewButton("New Capture Window", func() {
		pics.CapturePic()
		myApp.SendNotification(&fyne.Notification{
			Title:   "Wallpaper Tool",
			Content: "New Capture Window Opened",
		})
		getWindowsNum()
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

func getWindowsNum() {
	log.Println("Total Windows Opened: " + pics.GetLength(false))
}
