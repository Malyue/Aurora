package rabbitmq

import "github.com/streadway/amqp"

func NewRabbitMqSimple(queueName string, url string) (*RabbitMQ, error) {
	return NewRabbitMq(queueName, "", "", url)
}

func (r *RabbitMQ) PublishSimple(message string) error {
	_, err := r.channel.QueueDeclare(r.QueueName, false, false, false, false, nil)
	if err != nil {
		return err
	}

	// send msg to queue
	return r.channel.Publish(r.Exchange, r.QueueName, false, false, amqp.Publishing{ContentType: "text/plain", Body: []byte(message)})
}

func (r *RabbitMQ) ConsumeSimple() {

}
