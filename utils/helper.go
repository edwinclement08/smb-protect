package utils

import "fmt"

func ConnectAll(readOnly bool) {
	for _, shareMapping := range LoadedConfig.ShareMappings {
		roUser, roPass, rwUser, rwPass, err := LoadPasswords(shareMapping)
		if err != nil {
			fmt.Printf("Failed to load credentials for ShareMapping %s -> %s\n",
				shareMapping.SharePath, shareMapping.MountLocation)
		}

		user := rwUser
		pass := rwPass
		if readOnly {
			user = roUser
			pass = roPass
		}
		MountShare(shareMapping.MountLocation, shareMapping.SharePath, user, pass)
	}
}

func DisconnectAll() {
	for _, shareMapping := range LoadedConfig.ShareMappings {
		DisconnectShare(shareMapping.MountLocation)
	}
}
