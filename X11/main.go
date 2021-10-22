package main

import (
	"GUI/X11/pics"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"net/url"
)

func main() {
	go pics.MakeCache()

	var picWindows []fyne.Window
	myApp := app.New()
	mainWindow := myApp.NewWindow("Wallpaper Tool")
	mainWindow.SetMaster()
	mainWindow.Resize(fyne.NewSize(400, 100))

	captureBtn := widget.NewButton("Capture Picture", func() {
		go func() {
			picWindows = append(picWindows, pics.CapturePic(myApp))
		}()
	})

	closeBtn := widget.NewButton("Close All Pictures", func() {
		go pics.CloseAllWindows(picWindows)
	})
	bugURL, _ := url.Parse("https://github.com/sytgj7896321/GUI/issues/new")

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon(
			"Home",
			theme.HomeIcon(),
			container.New(layout.NewGridLayoutWithColumns(1), captureBtn, closeBtn),
		),
		container.NewTabItemWithIcon(
			"Favourite",
			theme.ListIcon(),
			container.New(layout.NewGridLayoutWithColumns(1), widget.NewLabel("TODO")),
		),
		container.NewTabItemWithIcon(
			"Settings",
			theme.SettingsIcon(),
			container.New(layout.NewGridLayoutWithColumns(1), widget.NewLabel("TODO")),
		),
		container.NewTabItemWithIcon(
			"Help",
			theme.HelpIcon(),
			container.New(layout.NewGridLayoutWithColumns(1), widget.NewHyperlink("Report a bug", bugURL))),
	)
	tabs.SetTabLocation(container.TabLocationTop)

	content := container.NewVBox(tabs)
	mainWindow.SetContent(content)
	mainWindow.Show()
	myApp.Run()
	tidyUp()
}

func tidyUp() {
	fmt.Println("Thank you for use")
	fmt.Println("Exited")
}
