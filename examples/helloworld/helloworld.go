package main

import (
	"github.com/zzl/goforms/forms"
	"github.com/zzl/goforms/framework/events"
)

func main() {

	form := forms.NewForm{
		Title:  "Hello world",
		Size:   forms.Size{Width: 600, Height: 400},
		Center: forms.CenterScreen,
	}.Create()

	forms.NewLabel{
		Text: "Name:",
		Pos:  forms.Pt(10, 10),
	}.Create()

	var txtName forms.Edit
	var btnGreet forms.Button

	txtName = forms.NewEdit{
		Text: "World",
		Pos:  forms.Pt(10, 30),
		Size: forms.Sz(60, 0),
		OnChange: func(ei *events.SimpleEventInfo) {
			btnGreet.SetEnabled(!txtName.IsEmpty())
		},
	}.Create()

	btnGreet = forms.NewButton{
		Text: "Greet",
		Pos:  forms.Pt(10, 64),
		Size: forms.Sz(60, 25),
		Action: func() {
			forms.Alert("Hello " + txtName.GetText() + "!")
		},
	}.Create()

	btnGreet.SetDefault()
	txtName.SelectAll()
	txtName.Focus()

	form.Show()
	form.Update()

	forms.MessageLoop()

}
