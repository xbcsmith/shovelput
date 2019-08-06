package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shopify/sarama"
	// TODO Deprecate cluster
	cluster "github.com/bsm/sarama-cluster"
)

// Consumer struct
type Consumer struct {
	consumer  *cluster.Consumer
	callbacks ConsumerCallbacks
}

// ConsumerCallbacks are callbacks for the consumer
type ConsumerCallbacks struct {
	OnDataReceived func(msg []byte)
	OnError        func(err error)
}

// NewConsumer returns a new Consumer
func NewConsumer(callbacks ConsumerCallbacks, brokerList []string, groupID string, topics []string) *Consumer {
	consumer := Consumer{callbacks: callbacks}

	config := cluster.NewConfig()
	config.ClientID = NewULID()
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	saramaConsumer, err := cluster.NewConsumer(brokerList, groupID, topics, config)
	if err != nil {
		panic(err)
	}
	consumer.consumer = saramaConsumer
	return &consumer

}

// Consume consumes messages from the bus
func (c *Consumer) Consume() {
	// Create signal channel
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	// Consume all channels, wait for signal to exit
	go func() {
		for {
			select {
			case msg, more := <-c.consumer.Messages():
				if more {
					if c.callbacks.OnDataReceived != nil {
						c.callbacks.OnDataReceived(msg.Value)
					}
					fmt.Fprintf(os.Stdout, "%s/%d/%d\t%s\n", msg.Topic, msg.Partition, msg.Offset, msg.Value)
					c.consumer.MarkOffset(msg, "")
				}
			case ntf, more := <-c.consumer.Notifications():
				if more {
					log.Printf("Rebalanced: %+v\n", ntf)
				}
			case err, more := <-c.consumer.Errors():
				if more {
					if c.callbacks.OnError != nil {
						c.callbacks.OnError(err)
					}
					//logger.Printf("Error: %s\n", err.Error())
				}
			case <-sigchan:
				return
			}
		}
	}()

}

// Close closes the consumer
func (c *Consumer) Close() {
	c.consumer.Close()
}
