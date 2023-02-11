package utils

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"syscall"
)

type ConnectionState struct {
	Connected bool
	Writable  bool
}

func DisconnectShare(driveLetter string) {
	cmd := exec.Command("C:/Windows/System32/net", "use", fmt.Sprintf("%s:", driveLetter), "/delete", "/y")

	var out bytes.Buffer
	var errorStream bytes.Buffer
	cmd.Stdin = strings.NewReader("")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Stdout = &out
	cmd.Stderr = &errorStream

	err := cmd.Run()

	if err != nil {
		if !strings.Contains(errorStream.String(), "The network connection could not be found.") {
			fmt.Printf("DisconnectShare:UnknownErrorMsg: %q\n", errorStream.String())
		}
	}
}

func CheckConnectedState(driveLetter string) ConnectionState {
	var connectionState ConnectionState
	cmd := exec.Command("C:/Windows/System32/net", "use", fmt.Sprintf("%s:", driveLetter))

	var out bytes.Buffer
	var errorStream bytes.Buffer
	cmd.Stdin = strings.NewReader("")
	cmd.Stdout = &out
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Stderr = &errorStream

	err := cmd.Run()

	if err != nil {
		if !strings.Contains(errorStream.String(), "The network connection could not be found.") {
			fmt.Printf("CheckConnnectedState:UnknownErrorMsg: %q\n", errorStream.String())
		}
		return connectionState
	}

	connectionState.Connected = true
	connectionState.Writable = IsWritable(fmt.Sprintf("%s:\\", driveLetter))
	return connectionState
}

func MountShare(driveLetter, sharePath, username, password string) {
	cmd := exec.Command("C:/Windows/System32/net", "use", fmt.Sprintf("%s:", driveLetter), sharePath, fmt.Sprintf("/USER:%s", username), "/persistent:no", fmt.Sprintf("%s", password))

	var out bytes.Buffer
	var errorStream bytes.Buffer
	cmd.Stdin = strings.NewReader("")
	cmd.Stdout = &out
	cmd.Stderr = &errorStream
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	err := cmd.Run()

	if err != nil {
		if !strings.Contains(errorStream.String(), "The local device name is already in use.") {
			fmt.Printf("MountShare:UnknownErrorMsg: %q\n", errorStream.String())
		}
		if !strings.Contains(errorStream.String(), "The network name cannot be found.") {
			fmt.Println("Network is unavailable")
			// TODO do something here
		}
	}
}
