package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io/ioutil"
	"log"
	"time"

	"github.com/Shopify/sarama"
)

// Producer struct for messaging
type Producer struct {
	asyncProducer sarama.AsyncProducer
	callbacks     ProducerCallbacks
	topic         string
}

// ProducerCallbacks struct
type ProducerCallbacks struct {
	OnError func(error)
}

// NewProducer creates a new producer
func NewProducer(callbacks ProducerCallbacks, brokerList []string, topic string, certFile *string, keyFile *string, caFile *string, verifySsl *bool) *Producer {
	producer := Producer{callbacks: callbacks, topic: topic}

	config := sarama.NewConfig()
	tlsConfig := createTLSConf(certFile, keyFile, caFile, verifySsl)
	if tlsConfig != nil {
		config.Net.TLS.Enable = true
		config.Net.TLS.Config = tlsConfig
	}
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Frequency = 500 * time.Millisecond

	saramaProducer, err := sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
		panic(err)
	}
	go func() {
		for err := range saramaProducer.Errors() {
			if producer.callbacks.OnError != nil {
				producer.callbacks.OnError(err)
			}
		}
	}()
	producer.asyncProducer = saramaProducer
	return &producer
}

// createTLSConf creates a new TLS configuration
func createTLSConf(certFile *string, keyFile *string, caFile *string, verifySsl *bool) (t *tls.Config) {
	if certFile != nil && keyFile != nil && caFile != nil && *certFile != "" && *keyFile != "" && *caFile != "" {
		cert, err := tls.LoadX509KeyPair(*certFile, *keyFile)
		if err != nil {
			log.Fatal(err)
		}

		caCert, err := ioutil.ReadFile(*caFile)
		if err != nil {
			log.Fatal(err)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		t = &tls.Config{
			Certificates:       []tls.Certificate{cert},
			RootCAs:            caCertPool,
			InsecureSkipVerify: *verifySsl,
		}
	}

	return t
}

// Produce produces messages
func (p *Producer) Produce(payload interface{}) {
	value := message{
		value: payload,
	}
	value.ensureEncoded()
	log.Println("producing: ", string(value.encoded))
	p.asyncProducer.Input() <- &sarama.ProducerMessage{
		Topic: p.topic,
		//Key:   sarama.StringEncoder(r.RemoteAddr),
		Value: &value,
	}
}

// message struct
type message struct {
	value   interface{}
	encoded []byte
	err     error
}

// ensureEncoded ensures messages are JSON
func (ale *message) ensureEncoded() {
	if ale.encoded == nil && ale.err == nil {
		ale.encoded, ale.err = json.Marshal(ale.value)
		if ale.err != nil {
			log.Println(ale.err)
		}
	}
}

// Length returns the length of a message
func (ale *message) Length() int {
	ale.ensureEncoded()
	return len(ale.encoded)
}

// Encode returns the encoded message
func (ale *message) Encode() ([]byte, error) {
	ale.ensureEncoded()
	return ale.encoded, ale.err
}

// Close closes the Producer
func (p *Producer) Close() error {
	log.Println("Producer.Close()")
	if err := p.asyncProducer.Close(); err != nil {
		return err
	}
	return nil
}
