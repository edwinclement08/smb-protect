package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/edwinclement08/smb-protect/utils"
)

var MainTree *widget.Tree
var CountAdds int = 0

type Tab struct {
	Title string
	View  func(w fyne.Window, f func()) fyne.CanvasObject
}

var (
	// Tabs defines the metadata for each tutorial
	Tabs = map[string]Tab{
		"welcome": {"Welcome", welcomeScreen},
	}

	// TabIndex defines how our tutorials should be laid out in the index tree
	TabIndex = map[string][]string{
		"": {"welcome"},
	}
)

func testScreen(_ fyne.Window, addToTutorials func()) fyne.CanvasObject {
	return container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle("Welcome to the test facility", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	))
}

func welcomeScreen(win fyne.Window, addToTutorials func()) fyne.CanvasObject {
	mountLocation := widget.NewEntry()
	mountDrive := widget.NewEntry()

	roUser := widget.NewEntry()
	rwUser := widget.NewEntry()
	roPass := widget.NewPasswordEntry()
	rwPass := widget.NewPasswordEntry()

	form := widget.NewForm(
		widget.NewFormItem("Mount Location", mountLocation),
		widget.NewFormItem("Drive Letter", mountDrive),
		widget.NewFormItem("User(Read-Only)", roUser),
		widget.NewFormItem("Pass(Read-Only)", roPass),
		widget.NewFormItem("User(Read-Write)", rwUser),
		widget.NewFormItem("Pass(Read-Write)", rwPass),
	)
	form.OnSubmit = func() {
		fmt.Println("Form submitted")
		utils.SaveShareMapping(roPass.Text, rwPass.Text, mountLocation.Text, roUser.Text, rwUser.Text)
	}
	form.OnCancel = func() {
		fmt.Println("Form cancelled")
		mountLocation.SetText("")
		mountDrive.SetText("")

		roUser.SetText("")
		rwUser.SetText("")
		roPass.SetText("")
		rwPass.SetText("")
	}
	// widget.NewButton("Add new leaf Node", func() {
	// 	addToTutorials()
	// }),

	return container.NewVBox(
		widget.NewLabelWithStyle("Add a new share", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		form,
	)
}

func addToTutorials() {
	if CountAdds == 0 {
		fmt.Println("Updating tutorials")
		Tabs["Added"] = Tab{"Test", testScreen}
		TabIndex[""] = append(TabIndex[""], "Added")
		MainTree.Refresh()
		MainTree.ScrollTo("binding")
		fmt.Println(Tabs)
		fmt.Println(TabIndex)
		CountAdds += 1
	}
}

func makeNav(setTutorial func(tutorial Tab, addtoTutorials func()), loadPrevious bool) *widget.Tree {
	tree := widget.NewTree(
		func(uid string) []string { // ChildUIDs
			return TabIndex[uid]
		},
		func(uid string) bool { // IsBranch
			children, ok := TabIndex[uid]

			return ok && len(children) > 0
		},
		func(branch bool) fyne.CanvasObject { // CreateNode
			return widget.NewLabel("Collection Widgets")
		},
		func(uid string, branch bool, obj fyne.CanvasObject) { // UpdateNode
			t, ok := Tabs[uid]
			if !ok {
				fyne.LogError("Missing tutorial panel: "+uid, nil)
				return
			}
			obj.(*widget.Label).SetText(t.Title)
			obj.(*widget.Label).TextStyle = fyne.TextStyle{}
		},
	)
	tree.OnSelected = func(uid string) {
		if t, ok := Tabs[uid]; ok {
			setTutorial(t, addToTutorials)
		}
	}
	tree.Select("welcome")
	MainTree = tree

	return tree
}
