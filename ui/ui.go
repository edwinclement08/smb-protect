package ui

import (
	"fmt"
	"log"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"github.com/danieljoos/wincred"
	"github.com/edwinclement08/smb-protect/utils"
)

func makeMenu(a fyne.App, w fyne.Window) *fyne.MainMenu {
	newItem := fyne.NewMenuItem("New", nil)
	checkedItem := fyne.NewMenuItem("Checked", nil)
	checkedItem.Checked = true
	disabledItem := fyne.NewMenuItem("Disabled", nil)
	disabledItem.Disabled = true
	otherItem := fyne.NewMenuItem("Other", nil)
	mailItem := fyne.NewMenuItem("Mail", func() { fmt.Println("Menu New->Other->Mail") })
	mailItem.Icon = theme.MailComposeIcon()
	otherItem.ChildMenu = fyne.NewMenu("",
		fyne.NewMenuItem("Project", func() { fmt.Println("Menu New->Other->Project") }),
		mailItem,
	)
	fileItem := fyne.NewMenuItem("File", func() { fmt.Println("Menu New->File") })
	fileItem.Icon = theme.FileIcon()
	dirItem := fyne.NewMenuItem("Directory", func() { fmt.Println("Menu New->Directory") })
	dirItem.Icon = theme.FolderIcon()
	newItem.ChildMenu = fyne.NewMenu("",
		fileItem,
		dirItem,
		otherItem,
	)

	openSettings := func() {
		fmt.Println("Settings pressed")
	}
	settingsItem := fyne.NewMenuItem("Settings", openSettings)
	settingsShortcut := &desktop.CustomShortcut{KeyName: fyne.KeyComma, Modifier: fyne.KeyModifierShortcutDefault}
	settingsItem.Shortcut = settingsShortcut
	w.Canvas().AddShortcut(settingsShortcut, func(shortcut fyne.Shortcut) {
		openSettings()
	})

	cutShortcut := &fyne.ShortcutCut{Clipboard: w.Clipboard()}
	cutItem := fyne.NewMenuItem("Cut", func() {
		shortcutFocused(cutShortcut, w)
	})
	cutItem.Shortcut = cutShortcut
	copyShortcut := &fyne.ShortcutCopy{Clipboard: w.Clipboard()}
	copyItem := fyne.NewMenuItem("Copy", func() {
		shortcutFocused(copyShortcut, w)
	})
	copyItem.Shortcut = copyShortcut
	pasteShortcut := &fyne.ShortcutPaste{Clipboard: w.Clipboard()}
	pasteItem := fyne.NewMenuItem("Paste", func() {
		shortcutFocused(pasteShortcut, w)
	})
	pasteItem.Shortcut = pasteShortcut
	performFind := func() { fmt.Println("Menu Find") }
	findItem := fyne.NewMenuItem("Find", performFind)
	findItem.Shortcut = &desktop.CustomShortcut{KeyName: fyne.KeyF, Modifier: fyne.KeyModifierShortcutDefault | fyne.KeyModifierAlt | fyne.KeyModifierShift | fyne.KeyModifierControl | fyne.KeyModifierSuper}
	w.Canvas().AddShortcut(findItem.Shortcut, func(shortcut fyne.Shortcut) {
		performFind()
	})

	helpMenu := fyne.NewMenu("Help",
		fyne.NewMenuItem("Documentation", func() {
			u, _ := url.Parse("https://developer.fyne.io")
			_ = a.OpenURL(u)
		}),
		fyne.NewMenuItem("Support", func() {
			u, _ := url.Parse("https://fyne.io/support/")
			_ = a.OpenURL(u)
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Sponsor", func() {
			u, _ := url.Parse("https://fyne.io/sponsor/")
			_ = a.OpenURL(u)
		}))

	// a quit item will be appended to our first (File) menu
	deleteCreds := fyne.NewMenuItem("Delete Creds", func() {
		fmt.Println("LIst of all creds")
		creds, err := utils.ListCreds()
		if err != nil {
			fmt.Println("Error in listing creds:", err)
			return
		}
		for _, cred := range creds {
			target := cred.TargetName
			genCred, err := wincred.GetGenericCredential(target)
			if err == nil {
				genCred.Delete()
			}
		}
		utils.LoadedConfig.ShareMappings = []utils.ShareMapping{}
		utils.SaveConfig()
		// TODO need to update UI
	})

	file := fyne.NewMenu("File", newItem, checkedItem, disabledItem, deleteCreds)
	device := fyne.CurrentDevice()
	if !device.IsMobile() && !device.IsBrowser() {
		file.Items = append(file.Items, fyne.NewMenuItemSeparator(), settingsItem)
	}
	main := fyne.NewMainMenu(
		file,
		fyne.NewMenu("Edit", cutItem, copyItem, pasteItem, fyne.NewMenuItemSeparator(), findItem),
		helpMenu,
	)
	checkedItem.Action = func() {
		checkedItem.Checked = !checkedItem.Checked
		main.Refresh()
	}
	return main
}

func shortcutFocused(s fyne.Shortcut, w fyne.Window) {
	switch sh := s.(type) {
	case *fyne.ShortcutCopy:
		sh.Clipboard = w.Clipboard()
	case *fyne.ShortcutCut:
		sh.Clipboard = w.Clipboard()
	case *fyne.ShortcutPaste:
		sh.Clipboard = w.Clipboard()
	}
	if focused, ok := w.Canvas().Focused().(fyne.Shortcutable); ok {
		focused.TypedShortcut(s)
	}
}

func setShortcuts(w fyne.Window) {
	ctrlTab := &desktop.CustomShortcut{KeyName: fyne.KeyTab, Modifier: fyne.KeyModifierControl}
	ctrlShiftTab := &desktop.CustomShortcut{KeyName: fyne.KeyTab, Modifier: fyne.KeyModifierControl | fyne.KeyModifierShift}

	log.Println("Setting up shortcuts")
	w.Canvas().AddShortcut(ctrlTab, func(shortcut fyne.Shortcut) {
		log.Println("We tapped Ctrl+Tab")
		SetSiblingNode(true)
	})
	w.Canvas().AddShortcut(ctrlShiftTab, func(shortcut fyne.Shortcut) {
		log.Println("We tapped Ctrl+shift+Tab")
		SetSiblingNode(false)
	})
	ctrlSpace := &desktop.CustomShortcut{KeyName: fyne.KeySpace, Modifier: fyne.KeyModifierControl}
	w.Canvas().AddShortcut(ctrlSpace, func(shortcut fyne.Shortcut) {
		println("We tapped Ctrl+Space")
	})

}

func SetupConfigWindow(app fyne.App, window fyne.Window) {
	window.Resize(fyne.Size{Width: 500, Height: 450})

	var _ fyne.Theme = (*CustomTheme)(nil)

	app.Settings().SetTheme(&CustomTheme{})

	setShortcuts(window)
	window.SetMainMenu(makeMenu(app, window))

	content := container.NewMax()
	setContent := func(t Tab, addPane func(utils.ShareMapping)) {
		content.Objects = []fyne.CanvasObject{t.View(window, addPane)}
		content.Refresh()
	}

	split := container.NewHSplit(makeNav(setContent), content)
	split.Offset = 0.2

	window.SetContent(split)
	window.Resize(fyne.NewSize(640, 460))
	window.ShowAndRun()
}
