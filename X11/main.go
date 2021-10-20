package main

import (
	"GUI/X11/pics"
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
	fmt.Println("Thank you for use")

	go pics.MakeCache()

	var windows []fyne.Window
	myApp := app.New()
	mainWindow := myApp.NewWindow("Wallpaper Tool")
	mainWindow.SetMaster()
	mainWindow.Resize(fyne.Size{
		Width:  300,
		Height: 100,
	})

	loginBtn := widget.NewButton("Capture Picture", func() {
		go func() {
			windows = append(windows, pics.CapturePic(myApp))
		}()
	})

	closeBtn := widget.NewButton("Close All Pictures", func() {
		go pics.CloseAllWindows(windows)
	})

	bugURL, _ := url.Parse("https://github.com/sytgj7896321/GUI/issues/new")

	content := container.NewVBox(loginBtn, closeBtn, widget.NewHyperlink("Report a bug", bugURL))

	mainWindow.SetContent(content)
	mainWindow.Show()
	myApp.Run()
	tidyUp()

}

func tidyUp() {
	fmt.Println("Exited")
}
