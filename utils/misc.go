package utils

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/google/uuid"
)

func CoalesceError(args ...error) error {
	for _, arg := range args {
		if arg != nil {
			return arg
		}
	}
	return nil
}

func IsWritable(location string) bool {
	fileName := fmt.Sprintf("smb-protect-%s-%s", uuid.New().String(), uuid.New().String())

	filePath := path.Join(location, fileName)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("Failed to open config file: %s\n", err.Error())
		errStr := err.Error()
		if !strings.Contains(errStr, "Access is denied.") {
			fmt.Println("Read-write check fails, with the following error: ", errStr)
		}
		return false
	}

	file.Close()
	os.Remove(filePath)

	return true
}
