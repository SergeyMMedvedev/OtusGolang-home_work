package producer

import (
	"fmt"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	logger  *slog.Logger
	channel *amqp.Channel
	amqpURI string
}

func NewProducer(logger *slog.Logger, amqpURI string) *Producer {
	return &Producer{
		logger:  logger,
		amqpURI: amqpURI,
	}
}

func (p *Producer) Init(
	exchangeName string,
	exchangeType string,
	durable bool,
	autoDelete bool,
	internal bool,
	noWait bool,
	TODO
) error {
	p.logger.Info(fmt.Sprintf("dialing %q", p.amqpURI))
	connection, err := amqp.Dial(p.amqpURI)
	if err != nil {
		return fmt.Errorf("Dial: %s", err)
	}
	defer connection.Close()

	p.logger.Info("got Connection, getting Channel")
	p.channel, err = connection.Channel()
	if err != nil {
		return fmt.Errorf("Channel: %s", err)
	}

	p.logger.Info(fmt.Sprintf("got Channel, declaring %q Exchange (%q)", exchangeType, exchange))
	if err := p.channel.ExchangeDeclare(
		exchangeName, // name
		exchangeType, // type
		durable,      // durable
		false,        // auto-deleted
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return fmt.Errorf("Exchange Declare: %s", err)
	}
	return nil
}
