package store

import (
	_message "Aurora/internal/apps/access-server/pkg/message"
	_kafka "Aurora/internal/pkg/kafka"
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type KafkaConsumer struct {
	consumer  sarama.Consumer
	cf        func(m *_message.ChatMessage)
	offlineCf func(m *_message.ChatMessage)
	channelCf func(m *_message.ChatMessage)
}

func NewKafkaConsumer(address []string) (*KafkaConsumer, error) {
	consumer, err := _kafka.NewKafkaConsumer(address)
	if err != nil {
		return nil, err
	}
	c := &KafkaConsumer{
		consumer: *consumer,
	}
	if err = c.run(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *KafkaConsumer) run() error {
	partitions, err := c.consumer.Partitions(KafkaChatMessageTopic)

	if err != nil {
		return err
	}

	for _, partition := range partitions {
		consumer, err := c.consumer.ConsumePartition(KafkaChatMessageTopic, partition, sarama.OffsetNewest)
		if err != nil {
			return err
		}

		go func(pc sarama.PartitionConsumer) {
			for m := range pc.Messages() {
				var cm = _message.ChatMessage{}
				err2 := _message.JsonCodec.Decode(m.Value, &cm)
				if err2 != nil {
					logrus.Error("message decode error %v", err2)
					continue
				}
				if c.cf != nil {
					c.cf(&cm)
				}
			}
		}(consumer)

		consumer2, err := c.consumer.ConsumePartition(KafkaChannelMessageTopic, partition, sarama.OffsetNewest)
		if err != nil {
			return err
		}

		go func(pc sarama.PartitionConsumer) {
			for m := range pc.Messages() {
				var cm = _message.ChatMessage{}
				err2 := _message.JsonCodec.Decode(m.Value, &cm)
				if err2 != nil {
					logrus.Error("message decode error %v", err2)
					continue
				}
				if c.channelCf != nil {
					c.channelCf(&cm)
				}
			}
		}(consumer2)

		consumer3, err := c.consumer.ConsumePartition(KafkaChatOfflineMessageTopic, partition, sarama.OffsetNewest)
		if err != nil {
			return err
		}

		go func(pc sarama.PartitionConsumer) {
			for m := range pc.Messages() {
				var cm = _message.ChatMessage{}
				err2 := _message.JsonCodec.Decode(m.Value, &cm)
				if err2 != nil {
					logrus.Error("message decode error %v", err2)
					continue
				}
				if c.offlineCf != nil {
					c.offlineCf(&cm)
				}
			}
		}(consumer3)
	}
	return nil
}

func (c *KafkaConsumer) Close() error {
	return c.consumer.Close()
}

func (c *KafkaConsumer) ConsumeChatMessage(cf func(m *_message.ChatMessage)) {
	c.cf = cf
}

func (c *KafkaConsumer) ConsumeChannelMessage(cf func(m *_message.ChatMessage)) {
	c.channelCf = cf
}

func (c *KafkaConsumer) ConsumeOfflineMessage(cf func(m *_message.ChatMessage)) {
	c.offlineCf = cf
}
