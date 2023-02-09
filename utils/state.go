package utils

import (
	"time"
)

type AppStateType struct {
	ConnectionStates map[string]ConnectionState // uuid -> connectionState
}

var AppState AppStateType

func UpdateConnectionStates() {
	if AppState.ConnectionStates == nil {
		AppState.ConnectionStates = make(map[string]ConnectionState)
	}

	shareMappings := LoadedConfig.ShareMappings
	if len(shareMappings) == 0 {
		AppState.ConnectionStates = map[string]ConnectionState{}
		return
	}

	for _, shareMapping := range shareMappings {
		state := CheckConnectedState(shareMapping.MountLocation)
		AppState.ConnectionStates[shareMapping.Uuid] = state
	}
}

func StateUpdateLoop() {
	for {
		UpdateConnectionStates()
		time.Sleep(time.Second * 5)
	}
}

func Check_AnyConnected() bool {
	connected := false
	for _, state := range AppState.ConnectionStates {
		if state.Connected {
			connected = true
			break
		}
	}
	return connected
}
