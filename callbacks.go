package main

// Callbacks struct for different callbacks
type Callbacks struct {
	OnNewConnection        func(clientID string)
	OnConnectionTerminated func(clientID string)
	OnDataReceived         func(clientID string, data []byte)
}
