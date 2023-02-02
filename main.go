package main

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2/app"

	"github.com/edwinclement08/smb-protect/ui"
	"github.com/edwinclement08/smb-protect/utils"
)

var AppConfig utils.ConfigType

func main() {
	a := app.NewWithID("com.edwinclement08.smb-protect")
	w := a.NewWindow("SMB Protect")

	// makeTray(a)
	// checkConnectedState("X")
	log.Println("Available Creds")
	utils.ListCred()
	log.Println("\n\nInitialize Config")
	// utils.InitConfig()

	config, err := utils.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load App Configuration: %s", err)
	}
	AppConfig = config
	for _, share := range config.ShareMappings {
		fmt.Printf("Adding share %s -> %s\n", share.SharePath, share.MountLocation)
		ui.AddPane(share)
	}

	// salt := utils.GenerateSeed()
	// pass := "pass23"
	// cipher, nonce := utils.Encrypt("Test encryption", pass, salt)
	// fmt.Printf("Encrypted text: %x\n", cipher)
	// plainText := utils.Decrypt(cipher, pass, nonce, salt)
	// fmt.Printf("Decrypted text: %s\n", plainText)

	ui.SetupConfigWindow(w)
}
