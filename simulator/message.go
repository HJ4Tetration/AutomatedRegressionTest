package main

type SwitchRegistrationResponse struct {
}

type SwitchRegistration struct {
	Serial string `json:"serial"`
	Crt    string `json:"crt"`
}

type SwitchMessage struct {
	Cmd      string `json:"cmd"`
	SwitchID string `json:"switchId"`
}

type ServerMessage struct {
	ResponseCode int    `json:"responseCode"`
	Cmd          string `json:"cmd"`
	// Data: ignore data field, use more specific structs if Data is needed
}
