package main

type Callbacks struct {
	OnNewConnection        func(clientID string)
	OnConnectionTerminated func(clientID string)
	OnDataReceived         func(clientID string, data []byte)
}
