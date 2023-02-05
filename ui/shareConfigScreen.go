package ui

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/edwinclement08/smb-protect/utils"
)

func CreateShareConfigScreen(shareMapping utils.ShareMapping) func(fyne.Window, func(utils.ShareMapping)) fyne.CanvasObject {
	roUser, roPass, rwUser, rwPass, err := utils.LoadPasswords(shareMapping)
	grayColor := color.RGBA{128, 128, 128, 255}
	greenColor := color.RGBA{97, 217, 124, 255}
	redColor := color.RGBA{223, 70, 85, 255}

	return func(_ fyne.Window, addPane func(utils.ShareMapping)) fyne.CanvasObject {
		mountLocation := widget.NewEntry()
		mountLocation.SetText(shareMapping.MountLocation)

		connectionStatus := &canvas.Text{
			Text:     "Checking",
			TextSize: theme.TextSize(),
			Color:    grayColor,
		}

		stateUpdate := make(chan bool)

		updateConnectionStatus := func() {
			connectionStatus.Color = grayColor
			connectionStatus.Text = "Checking"
			connectionStatus.Refresh()
			state := utils.CheckConnectedState(shareMapping.MountLocation)

			if state.Connected {
				connectionStatus.Color = greenColor
				readState := "RO"
				if state.Writable {
					readState = "RW"
				}
				connectionStatus.Text = fmt.Sprintf("Connected (%s)", readState)
			} else {
				connectionStatus.Color = redColor
				connectionStatus.Text = "Disconnected"

			}
			connectionStatus.Refresh()
			stateUpdate <- state.Connected
		}

		if err != nil {
			fmt.Println("Failed to load passwords for ", shareMapping.MountLocation)
			fmt.Println(err)
		}

		createConnectButton := func(readonly bool) *widget.Button {
			label := "Connect as Read/Write"
			user := rwUser
			pass := rwPass
			if readonly {
				label = "Connect as Read-Only"
				user = roUser
				pass = roPass
			}
			return widget.NewButton(label, func() {
				connectionStatus.Color = grayColor
				connectionStatus.Text = "Checking"
				connectionStatus.Refresh()
				utils.MountShare(shareMapping.MountLocation, shareMapping.SharePath, user, pass)
				updateConnectionStatus()
			})

		}

		connectRWButton := createConnectButton(false)
		connectROButton := createConnectButton(true)

		disconnectShareButton := widget.NewButton("Disconnect Share", func() {
			connectionStatus.Color = grayColor
			connectionStatus.Text = "Checking"
			connectionStatus.Refresh()
			utils.DisconnectShare(shareMapping.MountLocation)
			updateConnectionStatus()
		})

		checkStateButton := widget.NewButton("Check Connection", updateConnectionStatus)

		updateMountLocation := widget.NewButton("Update", func() {
			fmt.Println("Tapped the update button")
		})
		defer updateMountLocation.Disable()
		mountLocation.OnChanged = func(value string) {
			if value != shareMapping.MountLocation {
				updateMountLocation.Enable()
			} else {
				updateMountLocation.Disable()
			}
		}
		go updateConnectionStatus()

		biggerTextSize := theme.TextSize() + 2

		go func() {
			for {
				connected := <-stateUpdate

				if connected {
					connectROButton.Disable()
					connectRWButton.Disable()
					disconnectShareButton.Enable()
				} else {
					connectROButton.Enable()
					connectRWButton.Enable()
					disconnectShareButton.Disable()
				}
			}
		}()

		return NewBorderStyle(
			container.NewVBox(
				container.NewHBox(
					NewCustomBoldLabel("Share Map Configuration", color.Black, theme.TextSubHeadingSize()),
					layout.NewSpacer(),
					connectionStatus,
				),
				container.NewVBox(
					NewVSplitLayout(
						NewExpandedVBoxLayout(
							widget.NewLabelWithStyle("Share Path", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
							container.NewPadded(
								&canvas.Text{
									Text:     shareMapping.SharePath,
									TextSize: biggerTextSize,
									Color:    color.Black,
								},
							),
						),
						NewExpandedVBoxLayout(
							widget.NewLabelWithStyle("Mount Location", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
							NewVSplitLayout(
								mountLocation,
								updateMountLocation,
							),
						),
					),
					NewBorderStyle(
						NewExpandedVBoxLayout(
							NewCustomBoldLabel("Credentials", color.Black, biggerTextSize),
							container.NewGridWithColumns(4,
								widget.NewLabelWithStyle("Read-Only User ", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
								widget.NewLabel(shareMapping.ROUser),
								layout.NewSpacer(),
								widget.NewButton("Update", func() {
									fmt.Println("Update Read-only User button pressed")
								}),

								widget.NewLabelWithStyle("Read-Write User", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
								widget.NewLabel(shareMapping.RWUser),
								layout.NewSpacer(),
								widget.NewButton("Update", func() {
									fmt.Println("Update Read-write User button pressed")
								}),
							),
						),
					),
					NewBorderStyle(
						NewExpandedVBoxLayout(
							NewCustomBoldLabel("Advanced Config", color.Black, biggerTextSize),
							NewBorderStyle(
								widget.NewCheck("Automatically Mount at Login", func(b bool) { fmt.Println("Toggle Auto-mount") }),
								widget.NewRadioGroup([]string{"Read-Only", "Read-Write"}, func(val string) {
									fmt.Println(val, "is the mode of auto-mount")
								}),
							),
							widget.NewCheck("Open in explorer when mounted manually", func(b bool) { fmt.Println("Toggle Open in Explorer") }),
						),
					),

					layout.NewSpacer(),
					container.NewHBox(layout.NewSpacer(), checkStateButton),
					layout.NewSpacer(),
					container.NewHBox(connectROButton, connectRWButton, disconnectShareButton),
				)))
	}
}
