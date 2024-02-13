package kafka

import "github.com/IBM/sarama"

func NewKafkaConsumer(address []string) (*sarama.Consumer, error) {
	consumer, err := sarama.NewConsumer(address, sarama.NewConfig())
	if err != nil {
		return nil, err
	}

	return &consumer, nil
}
