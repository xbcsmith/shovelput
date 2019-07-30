package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/xbcsmith/shovelput/config"
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
	/*brokers   = flag.String("brokers", os.Getenv("KAFKA_PEERS"), "The Kafka brokers to connect to, as a comma separated list")
	producerTopic   = flag.String("producer_topic", "tcp_layer_messages", "The topic to produce messages to")
	consumerTopics   = flag.String("consumer_topic", "workers_layer_messages", "The topic to consume messages from")
	consumerGroupId   = flag.String("consumer_group_id", "", "consumer group id")
	verbose   = flag.Bool("verbose", false, "Turn on Sarama logging")
	certFile  = flag.String("certificate", "", "The optional certificate file for client authentication")
	keyFile   = flag.String("key", "", "The optional key file for client authentication")
	caFile    = flag.String("ca", "", "The optional certificate authority file for TLS client authentication")
	verifySsl = flag.Bool("verify", false, "Optional verify ssl certificates chain")*/
	configPath      = flag.String("config", "", "config file")
	consumerGroupId string
)

func main() {
	flag.Parse()
	if *configPath == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	config.InitConfig(*configPath)
	configuration := config.Get()
	/*if *brokers == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}*/

	if configuration.ConsumerGroupId == "" {
		consumerGroupId = NewULID()
	} else {
		consumerGroupId = configuration.ConsumerGroupId
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
	consumer = NewConsumer(consumerCallbacks, configuration.BrokersList, consumerGroupId, configuration.ConsumerTopics)
	consumer.Consume()

	signal_chan := make(chan os.Signal, 1)
	signal.Notify(signal_chan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGKILL)

	go func() {
		for {
			s := <-signal_chan
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
		server.SendDataByClientId(clientID, []byte("Pong"))
	}
	if producer != nil {
		var deviceRequest DeviceRequest
		err := json.Unmarshal(data, &deviceRequest)
		if err == nil {
			serverRequest := ServerRequest{
				DeviceRequest: deviceRequest,
				ServerId:      "1",
				ClientId:      clientID,
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
	if serverResponse.DeviceResponse.Action == "connect.response" && serverResponse.DeviceResponse.Status == "ok" && serverResponse.ClientId != "" {
		//attach the device id to our existing client
		err = server.SetDeviceIDToClient(serverResponse.ClientId, serverResponse.DeviceID)
		if err != nil {
			log.Println(err)
		}
	}
	toSend, err := json.Marshal(serverResponse.DeviceResponse)
	if err != nil {
		log.Println(err)
		return
	}
	if serverResponse.ClientId != "" {
		server.SendDataByClientId(serverResponse.ClientId, toSend)
	} else {
		if serverResponse.DeviceID != "" {
			server.SendDataByDeviceID(serverResponse.DeviceID, toSend)
		}
	}

}
