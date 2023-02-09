package ui

import (
	"fmt"
	"path"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"github.com/edwinclement08/smb-protect/utils"
)

func ShowHideConfigMenuItem(w fyne.Window, menu *fyne.Menu) *fyne.MenuItem {
	h := fyne.NewMenuItem("Hide Config", func() {})
	w.SetCloseIntercept(func() {
		w.Hide()
		h.Label = "Show Config"
		menu.Refresh()
	})

	h.Action = func() {
		if h.Label == "Show Config" {
			w.Show()
			h.Label = "Hide Config"
		} else if h.Label == "Hide Config" {
			w.Hide()
			h.Label = "Show Config"
		}
		menu.Refresh()
	}
	return h
}

func MakeTray(a fyne.App, w fyne.Window) {
	if desk, ok := a.(desktop.App); ok {
		icon, err := fyne.LoadResourceFromPath(path.Join(".", "resources", "smb-icon.jpg"))
		if err != nil {
			icon = theme.FileApplicationIcon()
		}
		var menu *fyne.Menu

		shares := fyne.NewMenuItem("Shares", func() {})

		shares.ChildMenu = fyne.NewMenu("Shares")

		result := []*fyne.MenuItem{}

		for _, sm := range utils.LoadedConfig.ShareMappings {
			state := utils.AppState.ConnectionStates[sm.Uuid]
			c := ""
			if state.Connected {
				c = "[RO]"
				if state.Writable {
					c = "[RW]"

				}
			}
			item := fyne.NewMenuItem(sm.SharePath+" "+c, func() {})
			item.Checked = state.Connected
			fmt.Println("Updating state", state)
			item.Action = func() {
				state := utils.AppState.ConnectionStates[sm.Uuid]
				if state.Connected {
					fmt.Println("Disconnecting")
					utils.DisconnectShare(sm.MountLocation)
					item.Checked = false
				} else {
					fmt.Println("Connecting")
					roUser, roPass, _, _, err := utils.LoadPasswords(sm)
					if err != nil {
						fmt.Println("failed to mount readonly", err.Error())
						return
					}
					utils.MountShare(sm.MountLocation, sm.SharePath, roUser, roPass)
					item.Label = sm.SharePath + " [RO]"
					item.Checked = true

				}

				shares.ChildMenu.Refresh()
				menu.Refresh()

			}

			result = append(result, item)
		}
		shares.ChildMenu.Items = result

		menu = fyne.NewMenu("smb-protect",
			ShowHideConfigMenuItem(w, menu),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("Connect All(RO)", func() {
				utils.ConnectAll(false)
			}),
			fyne.NewMenuItem("Connect All(RW)", func() {
				utils.ConnectAll(true)
			}),
			fyne.NewMenuItem("Disconnect All", func() {
				utils.DisconnectAll()
			}),
			fyne.NewMenuItemSeparator(),
			shares,
		)
		desk.SetSystemTrayIcon(icon)
		desk.SetSystemTrayMenu(menu)
	}
}
