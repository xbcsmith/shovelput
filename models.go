package main

type DeviceResponse struct {
	Action       string                 `json:"action"`
	ID           string                 `json:"id"`
	Data         map[string]interface{} `json:"data"`
	Status       string                 `json:"status"`
	ErrorCode    string                 `json:"errorCode"`
	ErrorMessage string                 `json:"errorMessage"`
}

type ServerResponse struct {
	DeviceResponse DeviceResponse `json:"deviceResponse"`
	ServerId       string         `json:"serverId"` // unique identifier for each server in case of having more than 1 server
	ClientId       string         `json:"clientId"`
	DeviceID       string         `json:"deviceID"`
}

type DeviceRequest struct {
	Action   string                 `json:"action"`
	DeviceID string                 `json:"deviceID"`
	ID       string                 `json:"id"`
	Data     map[string]interface{} `json:"data"`
}

type ServerRequest struct {
	DeviceRequest DeviceRequest `json:"deviceRequest"`
	ServerId      string        `json:"serverId"` // unique identifier for each server in case of having more than 1 server
	ClientId      string        `json:"clientId"`
}
