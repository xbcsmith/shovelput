package main

// DeviceResponse struct for device response
type DeviceResponse struct {
	Action       string                 `json:"action"`
	ID           string                 `json:"id"`
	Data         map[string]interface{} `json:"data"`
	Status       string                 `json:"status"`
	ErrorCode    string                 `json:"errorCode"`
	ErrorMessage string                 `json:"errorMessage"`
}

// ServerResponse struct for server response
type ServerResponse struct {
	DeviceResponse DeviceResponse `json:"deviceResponse"`
	ServerID       string         `json:"serverID"`
	ClientID       string         `json:"clientID"`
	DeviceID       string         `json:"deviceID"`
}

// DeviceRequest struct for device requests
type DeviceRequest struct {
	Action   string                 `json:"action"`
	DeviceID string                 `json:"deviceID"`
	ID       string                 `json:"id"`
	Data     map[string]interface{} `json:"data"`
}

// ServerRequest struct for server requests
type ServerRequest struct {
	DeviceRequest DeviceRequest `json:"deviceRequest"`
	ServerID      string        `json:"serverID"`
	ClientID      string        `json:"clientID"`
}
