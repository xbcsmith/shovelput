package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
)

// Server Struct
type Server struct {
	address     string // Address to open connection: localhost:9999
	connLock    sync.RWMutex
	connections map[string]*Client
	callbacks   Callbacks
	listener    net.Listener
}

// onConnectionEvent event processing for connection events
func (s *Server) onConnectionEvent(c *Client, eventType ConnectionEventType, e error) {
	switch eventType {
	case CONNECTION_EVENT_TYPE_NEW_CONNECTION:
		s.connLock.Lock()
		id := NewULID()
		c.ID = id
		s.connections[id] = c
		s.connLock.Unlock()
		//log.Println(eventType , " ,  id:", c.ID, " , ip: ", c.Conn().RemoteAddr().String())
		if s.callbacks.OnNewConnection != nil {
			s.callbacks.OnNewConnection(id)
		}
	case CONNECTION_EVENT_TYPE_CONNECTION_TERMINATED, CONNECTION_EVENT_TYPE_CONNECTION_GENERAL_ERROR:
		//log.Println(eventType , " ,  id:", c.ID, " , ip: ", c.Conn().RemoteAddr().String(), " , error: ", e.Error())
		s.connLock.Lock()
		delete(s.connections, c.ID)
		s.connLock.Unlock()
		if s.callbacks.OnConnectionTerminated != nil {
			s.callbacks.OnConnectionTerminated(c.ID)
		}
	}
}

// onDataEvent event processing for data events
func (s *Server) onDataEvent(c *Client, data []byte) {
	//log.Println("onDataEvent, ", c.Conn().RemoteAddr().String(), " data: " , string(data))
	if s.callbacks.OnDataReceived != nil {
		s.callbacks.OnDataReceived(c.ID, data)
	}
}

// Listen starts the client listener
func (s *Server) Listen() {
	var err error
	s.listener, err = net.Listen("tcp", s.address)
	if err != nil {
		log.Fatal("Error starting TCP Server.: ", err)
	}
	for {
		conn, _ := s.listener.Accept()
		client := NewClient(conn, s.onConnectionEvent, s.onDataEvent)
		s.onConnectionEvent(client, CONNECTION_EVENT_TYPE_NEW_CONNECTION, nil)
		go client.listen()

	}
}

// NewServer returns new Server instance
func NewServer(address string, callbacks Callbacks) *Server {
	log.Println("Creating Server with address", address)
	s := &Server{
		address:   address,
		callbacks: callbacks,
	}
	s.connections = make(map[string]*Client)
	return s
}

// SendDataByClientID sends data to a specific client id
func (s *Server) SendDataByClientID(clientID string, data []byte) error {
	if s.connections[clientID] == nil {
		return errors.New(fmt.Sprint("no connection with id ", clientID))
	}
	return s.connections[clientID].Send(data)
}

// SendDataByDeviceID sends data to specific device id
func (s *Server) SendDataByDeviceID(deviceID string, data []byte) error {
	for k := range s.connections {
		if s.connections[k].DeviceID == deviceID {
			return s.connections[k].Send(data)
		}
	}
	return errors.New(fmt.Sprint("no connection with deviceID ", deviceID))
}

// SetDeviceIDToClient sets the device id to the client id
func (s *Server) SetDeviceIDToClient(clientID string, deviceID string) error {
	if s.connections[clientID] != nil {
		s.connections[clientID].DeviceID = deviceID
		return nil
	}
	return errors.New(fmt.Sprint("no connection with id ", clientID))
}

// Close closes the server listener
func (s *Server) Close() {
	log.Println("Server.Close()")
	log.Println("s.connections length: ", len(s.connections))
	for k := range s.connections {
		fmt.Printf("key[%s]\n", k)
		s.connections[k].Close()
	}
	s.listener.Close()
}
