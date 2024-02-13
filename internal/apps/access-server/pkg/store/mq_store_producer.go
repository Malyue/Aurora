package store

import (
	_message "Aurora/internal/apps/access-server/pkg/message"
	_kafka "Aurora/internal/pkg/kafka"
	"encoding/json"
	"github.com/IBM/sarama"
	"time"
)

const (
	KafkaChatMessageTopic        = "getaway_chat_message"
	KafkaChatOfflineMessageTopic = "getaway_chat_offline_message"
	KafkaChannelMessageTopic     = "gateway_channel_message"
)

var _ MessageStore = &KafkaMessageStore{}

//var _ SubscriptionStore = &KafkaMessageStore{}

type KafkaMessageStore struct {
	producer sarama.AsyncProducer
}

type msg struct {
	data []byte
}

func NewKafkaProducer(address []string) (*KafkaMessageStore, error) {
	producer, err := _kafka.NewKafkaProducer(address)
	if err != nil {
		return nil, err
	}

	return &KafkaMessageStore{
		producer: *producer,
	}, nil
}

func (k *KafkaMessageStore) Close() error {
	return k.producer.Close()
}

func (m *msg) Encode() ([]byte, error) {
	return m.data, nil
}

func (m *msg) Length() int {
	return len(m.data)
}

func (k *KafkaMessageStore) StoreOffline(message *_message.ChatMessage) error {
	msgBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	cm := &sarama.ProducerMessage{
		Topic:     KafkaChatOfflineMessageTopic,
		Value:     &msg{data: msgBytes},
		Headers:   nil,
		Metadata:  nil,
		Offset:    0,
		Partition: 0,
		Timestamp: time.Now(),
	}
	k.producer.Input() <- cm
	return nil
}

func (k *KafkaMessageStore) StoreMessage(message *_message.ChatMessage) error {
	msgBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	cm := &sarama.ProducerMessage{
		Topic:     KafkaChatMessageTopic,
		Value:     &msg{data: msgBytes},
		Headers:   nil,
		Metadata:  nil,
		Offset:    0,
		Partition: 0,
		Timestamp: time.Now(),
	}

	k.producer.Input() <- cm
	return nil
}
