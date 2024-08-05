package consumer

import (
	"fmt"
	"log/slog"
	"os"

	c "github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/config"
	"github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/ringbuffer"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	logger  *slog.Logger
	conn    *amqp.Connection
	channel *amqp.Channel
	done    chan error
	tag     string
	amqpURI string
	Buf     *ringbuffer.RingBuffer
}

func NewConsumer(
	logger *slog.Logger,
	amqpURI string,
	consumerCfg c.ConsumerConf,
	exchange c.ExchangeConf,
	queueCfg c.QueueConf,
	bindingCfg c.BindingConf,
) (*Consumer, error) {
	c := &Consumer{
		conn:    nil,
		channel: nil,
		logger:  logger,
		tag:     consumerCfg.Tag,
		done:    make(chan error),
		amqpURI: amqpURI,
		Buf:     ringbuffer.NewRingBuffer(10),
	}
	var err error
	c.logger.Info("dialing " + amqpURI)
	c.conn, err = amqp.Dial(amqpURI)
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}
	go func() {
		c.logger.Info(fmt.Sprintf("closing: %s", <-c.conn.NotifyClose(make(chan *amqp.Error))))
	}()

	c.logger.Info("got Connection, getting Channel")
	c.channel, err = c.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("channel: %w", err)
	}
	c.logger.Info(fmt.Sprintf("got Channel, declaring Exchange (%q)", exchange.Name))
	if err = c.channel.ExchangeDeclare(
		exchange.Name,
		exchange.Type,
		exchange.Durable,
		exchange.AutoDelete,
		exchange.Internal,
		exchange.NoWait,
		exchange.Args,
	); err != nil {
		return nil, fmt.Errorf("exchange Declare: %w", err)
	}

	c.logger.Info(fmt.Sprintf("declared Exchange, declaring Queue %q", queueCfg.Name))

	queue, err := c.channel.QueueDeclare(
		queueCfg.Name,
		queueCfg.Durable,
		queueCfg.AutoDelete,
		queueCfg.Exclusive,
		queueCfg.NoWait,
		queueCfg.Args,
	)
	if err != nil {
		return nil, fmt.Errorf("queue Declare: %w", err)
	}
	c.logger.Info(
		fmt.Sprintf(
			"declared Queue (%q %d messages, %d consumers), binding to Exchange (key %q)",
			queue.Name, queue.Messages, queue.Consumers, queueCfg.Key,
		),
	)

	if err = c.channel.QueueBind(
		bindingCfg.QueueName,
		bindingCfg.Key,
		bindingCfg.Exchange,
		bindingCfg.NoWait,
		bindingCfg.Args,
	); err != nil {
		return nil, fmt.Errorf("queue Bind: %w", err)
	}
	c.logger.Info(
		fmt.Sprintf(
			"Queue bound to Exchange, starting Consume (consumer tag %q)", c.tag,
		),
	)
	deliveries, err := c.channel.Consume(
		queue.Name,
		consumerCfg.Tag,
		consumerCfg.NoAck,
		consumerCfg.Exclusive,
		consumerCfg.NoLocal,
		consumerCfg.NoWait,
		consumerCfg.Args,
	)
	if err != nil {
		return nil, fmt.Errorf("queue Consume: %w", err)
	}

	go handle(deliveries, c.done, c.Buf)
	return c, nil
}

func (c *Consumer) Shutdown() error {
	if err := c.channel.Cancel(c.tag, true); err != nil {
		return fmt.Errorf("Consumer cancel failed: %w", err)
	}

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %w", err)
	}

	defer c.logger.Info("AMQP shutdown OK")

	return <-c.done
}

func handle(deliveries <-chan amqp.Delivery, done chan error, buf *ringbuffer.RingBuffer) {
	deliveredMsgs := "/tmp/delivered.txt"

	for d := range deliveries {
		msg := fmt.Sprintf(
			"got %dB delivery: [%v] %q",
			len(d.Body),
			d.DeliveryTag,
			d.Body,
		)
		slog.Info(msg)
		buf.Enqueue(d.Body)
		d.Ack(false)
		file, err := os.OpenFile(deliveredMsgs, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			slog.Error("open file error:" + err.Error())
		}
		_, err = file.WriteString(string(d.Body) + "\n")
		file.Close()
		if err != nil {
			slog.Error("write to file error:" + err.Error())
		}
	}
	slog.Info("handle: deliveries channel closed")
	done <- nil
}
