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
	arguments map[string]interface{},
	reliable bool,
) error {
	p.logger.Info(fmt.Sprintf("dialing %q", p.amqpURI))
	connection, err := amqp.Dial(p.amqpURI)
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}
	defer connection.Close()

	p.logger.Info("got Connection, getting Channel")
	p.channel, err = connection.Channel()
	if err != nil {
		return fmt.Errorf("channel: %w", err)
	}

	p.logger.Info(fmt.Sprintf("got Channel, declaring %q Exchange (%q)", exchangeType, exchangeName))
	if err := p.channel.ExchangeDeclare(
		exchangeName, // name
		exchangeType, // type
		durable,      // durable
		autoDelete,   // auto-deleted
		internal,     // internal
		noWait,       // noWait
		arguments,    // arguments
	); err != nil {
		return fmt.Errorf("exchange Declare: %w", err)
	}

	if reliable {
		p.logger.Info("enabling publishing confirms.")
		if err := p.channel.Confirm(false); err != nil {
			return fmt.Errorf("channel could not be put into confirm mode: %s", err)
		}

		confirms := p.channel.NotifyPublish(make(chan amqp.Confirmation, 1))

		defer p.confirmOne(confirms)
	}
	return nil
}

func (p *Producer) confirmOne(confirms <-chan amqp.Confirmation) {
	p.logger.Info("waiting for confirmation of one publishing")

	if confirmed := <-confirms; confirmed.Ack {
		p.logger.Info(fmt.Sprintf("confirmed delivery with delivery tag: %d", confirmed.DeliveryTag))
	} else {
		p.logger.Info(fmt.Sprintf("failed delivery of delivery tag: %d", confirmed.DeliveryTag))
	}
}

func (p *Producer) Publish(
	exchange,
	key string,
	mandatory, immediate bool,
	body string,
) error {
	if err := p.channel.Publish(
		exchange,  // Publish to an exchange
		key,       // Routing to a queue
		mandatory, // MANDATORY flag
		immediate, // IMMEDIATE flag
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            []byte(body),
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
			// a bunch of application/implementation-specific fields
		},
	); err != nil {
		return fmt.Errorf("error publishing: %w", err)
	}
	return nil
}
