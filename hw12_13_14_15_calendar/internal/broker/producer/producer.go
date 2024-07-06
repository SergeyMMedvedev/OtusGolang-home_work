package producer

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	brSchemas "github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/broker/schemas"
	c "github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/config"
	s "github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/storage"
	"github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/storage/schemas"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	logger   *slog.Logger
	channel  *amqp.Channel
	conn     *amqp.Connection
	amqpURI  string
	eventsCh chan []schemas.Event
}

func NewProducer(
	logger *slog.Logger,
	amqpURI string,
	exchangeCfg c.ExchangeConf,
) (*Producer, error) {
	p := &Producer{
		logger:   logger,
		amqpURI:  amqpURI,
		eventsCh: make(chan []schemas.Event),
	}

	p.logger.Info(fmt.Sprintf("dialing %q", p.amqpURI))
	var err error
	p.conn, err = amqp.Dial(p.amqpURI)
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}

	p.logger.Info("got Connection, getting Channel")
	p.channel, err = p.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("channel: %w", err)
	}

	p.logger.Info(fmt.Sprintf("got Channel, declaring %q Exchange (%q)", exchangeCfg.Type, exchangeCfg.Name))
	if err := p.channel.ExchangeDeclare(
		exchangeCfg.Name,
		exchangeCfg.Type,
		exchangeCfg.Durable,
		exchangeCfg.AutoDelete,
		exchangeCfg.Internal,
		exchangeCfg.NoWait,
		exchangeCfg.Args,
	); err != nil {
		return nil, fmt.Errorf("exchange Declare: %w", err)
	}

	if exchangeCfg.Reliable {
		p.logger.Info("enabling publishing confirms.")
		if err := p.channel.Confirm(false); err != nil {
			return nil, fmt.Errorf("channel could not be put into confirm mode: %w", err)
		}

		confirms := p.channel.NotifyPublish(make(chan amqp.Confirmation, 1))

		defer p.confirmOne(confirms)
	}

	return p, nil
}

func (p *Producer) Shutdown() error {
	if err := p.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %w", err)
	}
	defer p.logger.Info("AMQP shutdown OK")
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
			DeliveryMode:    amqp.Transient,
			Priority:        0,
		},
	); err != nil {
		return fmt.Errorf("error publishing: %w", err)
	}
	return nil
}

func (p *Producer) PublishEventsForNotification(
	wg *sync.WaitGroup,
	exchangeCfg c.ExchangeConf,
) {
	defer wg.Done()
	for events := range p.eventsCh {
		for _, event := range events {
			p.logger.Info("send notification for event: " + event.String())
			err := p.Publish(
				exchangeCfg.Name,
				exchangeCfg.Key,
				false,
				false,
				brSchemas.Notification{
					EventID:    event.ID,
					EventTitle: event.Title,
					EventDate:  event.Date,
					UserID:     event.UserID,
				}.String(),
			)
			if err != nil {
				p.logger.Error("failed to publish event " + event.ID + ", error: " + err.Error())
				continue
			}
		}
	}
}

func (p *Producer) SearchEventsForNotification(
	ctx context.Context,
	ticker *time.Ticker,
	storage s.Storage,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	for range ticker.C {
		p.logger.Info("call SearchEventsForNotification")
		events, err := storage.ListEventsForNotification(ctx)
		if err != nil {
			p.logger.Error("failed to list events: " + err.Error())
			continue
		}
		p.eventsCh <- events
	}
}

func (p *Producer) SearchAndRemoveOldEvents(
	ctx context.Context,
	ticker *time.Ticker,
	storage s.Storage,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	for range ticker.C {
		p.logger.Info("call SearchAndRemoveOldEvents")
		events, err := storage.ListLastYearEvents(ctx)
		if err != nil {
			p.logger.Error("failed to list events: " + err.Error())
			continue
		}
		for _, event := range events {
			p.logger.Info("Remove old event: " + event.ID)
			storage.DeleteEvent(ctx, event.ID)
		}
	}
}
