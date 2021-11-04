package main

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("List Data")

	data := binding.BindStringList(
		&[]string{"Item 1", "Item 2", "Item 3"},
	)

	list := widget.NewListWithData(data,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})

	add := widget.NewButton("Append", func() {
		val := fmt.Sprintf("Item %d", data.Length()+1)
		_ = data.Append(val)
	})
	tBtn := widget.NewButton("test", func() {
		err := errors.New("2333")
		win := fyne.CurrentApp().NewWindow("Error")
		win.Resize(fyne.NewSize(300, 200))
		win.SetContent(widget.NewTextGridFromString(err.Error()))
		win.Show()
	})

	myWindow.SetContent(container.NewBorder(tBtn, add, nil, nil, list))
	myWindow.ShowAndRun()
}
