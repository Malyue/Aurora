package kafka

import (
	"github.com/IBM/sarama"
)

func NewKafkaProducer(address []string) (*sarama.AsyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner

	producer, err := sarama.NewAsyncProducer(address, config)
	if err != nil {
		return nil, err
	}

	return &producer, nil
}
