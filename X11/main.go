package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
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
	myWin.Resize(fyne.Size{
		Width:  300,
		Height: 300,
	})

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("input name")
	//nameEntry.OnChanged = func(content string) {
	//	fmt.Println("name:", nameEntry.Text, "entered")
	//}

	passEntry := widget.NewPasswordEntry()
	passEntry.SetPlaceHolder("input password")

	nameBox := container.NewVBox(widget.NewLabel("Name"), nameEntry)

	passwordBox := container.NewVBox(widget.NewLabel("Password"), passEntry)

	loginBtn := widget.NewButton("Login", func() {
		fmt.Println("name:", nameEntry.Text, "password:", passEntry.Text, "login in")
	})

	multiEntry := widget.NewEntry()
	multiEntry.SetPlaceHolder("please enter\nyour description")
	multiEntry.MultiLine = true

	content := container.NewVBox(nameBox, passwordBox, loginBtn)

	myWin.SetContent(content)
	myWin.ShowAndRun()
}
