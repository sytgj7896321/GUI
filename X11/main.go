package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"net/url"
)

type myTheme struct{}

func (m myTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m myTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	if name == theme.IconNameHome {
		homeBytes := []byte{1, 2, 3}
		fyne.NewStaticResource("myHome", homeBytes)
	}
	return theme.DefaultTheme().Icon(name)
}

func (m myTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

func (m myTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameBackground {
		if variant == theme.VariantLight {
			return color.White
		}
		return color.Black
	}

	return theme.DefaultTheme().Color(name, variant)
}

func main() {
	var _ fyne.Theme = (*myTheme)(nil)
	myApp := app.New()
	myWin := myApp.NewWindow("Demo")
	myWin.SetMaster()
	myWin.Resize(fyne.Size{
		Width:  300,
		Height: 300,
	})

	//Name Box
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("input name")
	nameBox := container.NewVBox(widget.NewLabel("Name"), nameEntry)

	//Password Box
	passEntry := widget.NewPasswordEntry()
	passEntry.SetPlaceHolder("input password")
	passwordBox := container.NewVBox(widget.NewLabel("Password"), passEntry)

	//Description Box
	multiEntry := widget.NewEntry()
	multiEntry.SetPlaceHolder("please enter\nyour description")
	multiEntry.MultiLine = true

	//Login Button
	loginBtn := widget.NewButton("Login", func() {
		go showLogin(myApp, nameEntry.Text, passEntry.Text, multiEntry.Text)
	})

	//URL Button
	bugURL, _ := url.Parse("https://github.com/sytgj7896321/GUI/issues/new")

	content := container.NewVBox(nameBox, passwordBox, multiEntry, loginBtn, widget.NewHyperlink("Report a bug", bugURL))

	myWin.SetContent(content)
	myWin.Show()
	myApp.Run()
	tidyUp()

}

func tidyUp() {
	fmt.Println("Exited")
}

func showLogin(app fyne.App, name, password, description string) {
	win := app.NewWindow("Login Window")
	win.SetContent(widget.NewLabel("name: " + name + " password: " + password + " login in\nDescription: " + description))
	win.Resize(fyne.NewSize(200, 200))
	win.Show()
}
