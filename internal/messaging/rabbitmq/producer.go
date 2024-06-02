package rabbitmq

import (
	"github.com/streadway/amqp"
)

type Producer struct {
	conn *amqp.Connection
}

func NewProducer(conn *amqp.Connection) *Producer {
	return &Producer{conn: conn}
}

func (p *Producer) Publish(message []byte, routingKey string) error {
	channel, err := p.conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	err = channel.Publish(
		"",         // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		})
	if err != nil {
		return err
	}
	return nil
}
