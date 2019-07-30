package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var server *Server
var producer *Producer
var consumer *Consumer

var (
	brokers         = flag.String("brokers", os.Getenv("KAFKA_PEERS"), "The Kafka brokers to connect to, as a comma separated list")
	producerTopic   = flag.String("producer_topic", "tcp_layer_messages", "The topic to produce messages to")
	consumerTopics  = flag.String("consumer_topic", "workers_layer_messages", "The topic to consume messages from")
	consumerGroupID = flag.String("consumer_group_id", "", "consumer group id")
	verbose         = flag.Bool("verbose", false, "Turn on Sarama logging")
	certFile        = flag.String("certificate", "", "The optional certificate file for client authentication")
	keyFile         = flag.String("key", "", "The optional key file for client authentication")
	caFile          = flag.String("ca", "", "The optional certificate authority file for TLS client authentication")
	verifySsl       = flag.Bool("verify", false, "Optional verify ssl certificates chain")
)

func main() {
	flag.Parse()
	if brokers == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	configuration := NewConfig(brokers, producerTopic, consumerTopics, consumerGroupID)
	if configuration.ConsumerGroupID == "" {
		consumerGroupID = NewULID()
	}

	log.Printf("Kafka brokers: %s", strings.Join(configuration.BrokersList, ", "))
	callbacks := Callbacks{
		OnDataReceived:         onDataReceived,
		OnConnectionTerminated: onConnectionTerminated,
		OnNewConnection:        onNewConnection,
	}
	server = NewServer(":3000", callbacks)
	producerCallbacks := ProducerCallbacks{
		OnError: onProducerError,
	}
	f := false
	producer = NewProducer(producerCallbacks, configuration.BrokersList, configuration.ProducerTopic, nil, nil, nil, &f)

	consumerCallbacks := ConsumerCallbacks{
		OnDataReceived: onDataConsumed,
		OnError:        onConsumerError,
	}
	consumer = NewConsumer(consumerCallbacks, configuration.BrokersList, consumerGroupID, configuration.ConsumerTopics)
	consumer.Consume()

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGKILL)

	go func() {
		for {
			s := <-signalChannel
			switch s {
			case syscall.SIGINT:
				fmt.Println("syscall.SIGINT")
				cleanup()
				// kill -SIGTERM XXXX
			case syscall.SIGTERM:
				fmt.Println("syscall.SIGTERM")
				cleanup()
				// kill -SIGQUIT XXXX
			case syscall.SIGQUIT:
				fmt.Println("syscall.SIGQUIT")
				cleanup()
			case syscall.SIGKILL:
				fmt.Println("syscall.SIGKILL")
				cleanup()
			default:
				fmt.Println("Unknown signal.")
			}
		}
	}()

	go func() {
		http.HandleFunc("/", handler)
		http.ListenAndServe(":8080", nil)
	}()

	server.Listen()

}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Server is up and running!")
}

func cleanup() {
	server.Close()
	producer.Close()
	consumer.Close()
	os.Exit(0)
}

func onNewConnection(clientID string) {
	log.Println("onNewConnection, id: ", clientID)
}

func onConnectionTerminated(clientID string) {
	log.Println("onConnectionTerminated, id: ", clientID)
}

/**
Called when data is received from a TCP client, will generate a message to the message broker
*/
func onDataReceived(clientID string, data []byte) {
	log.Println("onDataReceived, id: ", clientID, ", data: ", string(data))
	if string(data) == "Ping" {
		log.Println("sending Pong")
		//answer with pong
		server.SendDataByClientID(clientID, []byte("Pong"))
	}
	if producer != nil {
		var deviceRequest DeviceRequest
		err := json.Unmarshal(data, &deviceRequest)
		if err == nil {
			serverRequest := ServerRequest{
				DeviceRequest: deviceRequest,
				ServerID:      "1",
				ClientID:      clientID,
			}
			producer.Produce(serverRequest)
		} else {
			log.Println(err)
		}

	}

}

func onProducerError(err error) {
	log.Println("onProducerError: ", err)
}

func onConsumerError(err error) {
	log.Println("onConsumerError: ", err)
}

func onDataConsumed(data []byte) {
	log.Println("onDataConsumed: ", string(data))
	var serverResponse ServerResponse
	err := json.Unmarshal(data, &serverResponse)
	if err != nil {
		log.Println(err)
		return
	}
	if serverResponse.DeviceResponse.Action == "connect.response" && serverResponse.DeviceResponse.Status == "ok" && serverResponse.ClientID != "" {
		//attach the device id to our existing client
		err = server.SetDeviceIDToClient(serverResponse.ClientID, serverResponse.DeviceID)
		if err != nil {
			log.Println(err)
		}
	}
	toSend, err := json.Marshal(serverResponse.DeviceResponse)
	if err != nil {
		log.Println(err)
		return
	}
	if serverResponse.ClientID != "" {
		server.SendDataByClientID(serverResponse.ClientID, toSend)
	} else {
		if serverResponse.DeviceID != "" {
			server.SendDataByDeviceID(serverResponse.DeviceID, toSend)
		}
	}

}
