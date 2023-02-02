package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/edwinclement08/smb-protect/utils"
)

var MainTree *widget.Tree

type Tab struct {
	Title string
	View  func(w fyne.Window, f func(utils.ShareMapping)) fyne.CanvasObject
}

var (
	// Tabs defines the metadata for each tutorial
	Tabs = map[string]Tab{
		"addShare": {"Add New Share", addShareScreen},
	}

	// TabIndex defines how our tutorials should be laid out in the index tree
	TabIndex = map[string][]string{
		"": {"addShare"},
	}
)

func AddPane(shareMapping utils.ShareMapping) {
	Tabs[shareMapping.Uuid] = Tab{
		fmt.Sprintf("%s -> %s:", shareMapping.SharePath, shareMapping.MountLocation),
		createShareConfigScreen(shareMapping),
	}
	TabIndex[""] = append(TabIndex[""], shareMapping.Uuid)
	if MainTree != nil {
		fmt.Println("Refreshing the Tree")
		MainTree.Refresh()
	}
}

func createShareConfigScreen(shareMapping utils.ShareMapping) func(fyne.Window, func(utils.ShareMapping)) fyne.CanvasObject {
	return func(_ fyne.Window, addPane func(utils.ShareMapping)) fyne.CanvasObject {
		// return container.NewCenter(container.NewVBox(
		// 	widget.NewLabelWithStyle("Welcome to the test facility", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		// ))

		return container.NewCenter(container.NewVBox(
			widget.NewLabel("Share Map Configuration"),
			widget.NewLabel(fmt.Sprintf("%s -> %s:", shareMapping.SharePath, shareMapping.MountLocation)),
		))
	}
}

func addShareScreen(win fyne.Window, addPane func(utils.ShareMapping)) fyne.CanvasObject {
	sharePath := widget.NewEntry()
	mountLocation := widget.NewEntry()

	roUser := widget.NewEntry()
	rwUser := widget.NewEntry()
	roPass := widget.NewPasswordEntry()
	rwPass := widget.NewPasswordEntry()

	form := widget.NewForm(
		widget.NewFormItem("Share Path", sharePath),
		widget.NewFormItem("Drive Letter", mountLocation),
		widget.NewFormItem("User(Read-Only)", roUser),
		widget.NewFormItem("Pass(Read-Only)", roPass),
		widget.NewFormItem("User(Read-Write)", rwUser),
		widget.NewFormItem("Pass(Read-Write)", rwPass),
	)
	form.OnSubmit = func() {
		fmt.Println("Form submitted")
		if sharePath.Text == "" ||
			roPass.Text == "" ||
			rwPass.Text == "" ||
			mountLocation.Text == "" ||
			roUser.Text == "" ||
			rwUser.Text == "" {
			fmt.Println("Field invalid")
			return
		}

		shareMapping := utils.SaveShareMapping(sharePath.Text, roUser.Text, roPass.Text, rwUser.Text, rwPass.Text, mountLocation.Text)
		addPane(shareMapping)

		sharePath.SetText("")
		mountLocation.SetText("")

		roUser.SetText("")
		rwUser.SetText("")
		roPass.SetText("")
		rwPass.SetText("")
	}
	form.OnCancel = func() {
		fmt.Println("Form cancelled")
		sharePath.SetText("")
		mountLocation.SetText("")

		roUser.SetText("")
		rwUser.SetText("")
		roPass.SetText("")
		rwPass.SetText("")
	}

	return container.NewVBox(
		widget.NewLabelWithStyle("Add a new share", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		form,
	)
}

func makeNav(setContent func(tab Tab, addPane func(utils.ShareMapping))) *widget.Tree {
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
			setContent(t, AddPane)
		}
	}
	MainTree = tree
	tree.Select("addShare")

	return tree
}
