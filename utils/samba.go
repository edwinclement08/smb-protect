package utils

import (
	"bytes"
	"fmt"
	"os/exec"
)

func DisconnectShare(driveLetter string) string {
	cmd := exec.Command("C:/Windows/System32/net", "use", fmt.Sprintf("%s:", driveLetter), "/delete", "/y")

	var out bytes.Buffer
	var errorStream bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errorStream

	err := cmd.Run()

	if err != nil {
		fmt.Printf("translated phrase: %q\n", out.String())
		fmt.Printf("translated phrase: %q\n", errorStream.String())
		return errorStream.String()
	}

	fmt.Printf("translated phrase: %q\n", out.String())
	return out.String()
}

func CheckConnectedState(driveLetter string) bool {
	cmd := exec.Command("C:/Windows/System32/net", "use", fmt.Sprintf("%s:", driveLetter))

	var out bytes.Buffer
	var errorStream bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errorStream

	err := cmd.Run()

	if err != nil {
		fmt.Printf("error")
		fmt.Printf("translated phrase: %q\n", out.String())
		fmt.Printf("translated phrase: %q\n", errorStream.String())
		return false
	}

	fmt.Printf("no error phrase: %q\n", out.String())
	return true
}

func MountShare(driveLetter, sharePath, username, password string) string {
	cmd := exec.Command("C:/Windows/System32/net", "use", fmt.Sprintf("%s:", driveLetter), sharePath, fmt.Sprintf("/USER:%s", username), "/persistent:no", fmt.Sprintf("%s", password))

	var out bytes.Buffer
	var errorStream bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errorStream

	err := cmd.Run()

	if err != nil {
		fmt.Printf("translated phrase: %q\n", out.String())
		fmt.Printf("translated phrase: %q\n", errorStream.String())
		return errorStream.String()
	}

	fmt.Printf("translated phrase: %q\n", out.String())
	return out.String()
}
