package main

import (
	"log"

	"fyne.io/fyne/v2/app"

	"github.com/edwinclement08/smb-protect/ui"
	"github.com/edwinclement08/smb-protect/utils"
)

var AppConfig utils.ConfigType

func main() {
	a := app.NewWithID("com.edwinclement08.smb-protect")
	w := a.NewWindow("SMB Protect")

	config, err := utils.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load App Configuration: %s", err)
	}
	AppConfig = config
	for _, share := range config.ShareMappings {
		ui.AddPane(share)
	}

	go utils.StateUpdateLoop()
	ui.MakeTray(a, w)
	ui.SetupConfigWindow(a, w)
}
