package main

// Config struct
type Config struct {
	BrokersList     []string `json:"brokers_list"`
	ProducerTopic   string   `json:"producer_topic"`
	ConsumerTopics  []string `json:"consumer_topics"`
	ConsumerGroupID string   `json:"consumer_group_id"`
}

// NewConfig returns a new config object
func NewConfig(brokers []string, producerTopic string, consumerTopics []string, consumerGroupID string) (*Config, error) {
	return &Config{
		BrokersList:     brokers,
		ProducerTopic:   producerTopic,
		ConsumerTopics:  consumerTopics,
		ConsumerGroupID: consumerGroupID,
	}, nil
}
