package ui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"github.com/edwinclement08/smb-protect/utils"
)

func makeTray(a fyne.App) {
	if desk, ok := a.(desktop.App); ok {
		h := fyne.NewMenuItem("Hello", func() {})
		h.Icon = theme.HomeIcon()
		menu := fyne.NewMenu("Hello World", h)
		h.Action = func() {
			log.Println("System tray menu tapped")
			h.Label = "Welcome"
			menu.Refresh()
		}
		desk.SetSystemTrayMenu(menu)
	}
}

func SetupConfigWindow(window fyne.Window) {
	window.Resize(fyne.Size{Width: 350, Height: 250})

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

func SetupConnectionStatus() {
	// connectionStatus := widget.NewLabel("Connection State: Checking")

	// commandOutputView := widget.NewTextGrid()
	// connectShareButton := widget.NewButton("Connect Share", func() {
	// 	log.Println("Button Pressed")
	// 	commandOutputView.SetText("Connecting...")
	// 	output := utils.MountShare("X", "\\\\192.168.0.192\\Edwin", "edwin-rw", "password")
	// 	commandOutputView.SetText(fmt.Sprint("Connected\n", output))
	// })

	// disconnectShareButton := widget.NewButton("Disconnect Share", func() {
	// 	commandOutputView.SetText("Disconnecting...")
	// 	output := utils.DisconnectShare("X")
	// 	commandOutputView.SetText(fmt.Sprint("Disconnected\n", output))
	// })
	// checkStateButton := widget.NewButton("Check Connection", func() {
	// 	commandOutputView.SetText("")
	// 	connectionStatus.SetText("Connection State: Checking")
	// 	connected := utils.CheckConnectedState("X")
	// 	updatedText := "Connection State: offline"
	// 	if connected {
	// 		updatedText = "Connection State: online"
	// 	}
	// 	connectionStatus.SetText(updatedText)
	// })
	// horizontalLayout := container.NewHBox(checkStateButton, connectShareButton, disconnectShareButton)
	// mainLayout := container.NewVBox(connectionStatus, horizontalLayout, commandOutputView)

	// content := container.New(layout.NewGridLayout(2), text1, text2)
	// window.SetContent(mainLayout)

	// window.SetCloseIntercept(func() {
	// 	window.Hide()
	// })

	// window.ShowAndRun()

}
