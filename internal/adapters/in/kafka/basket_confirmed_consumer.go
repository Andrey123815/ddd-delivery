package kafkain

import (
	"context"

	createOrder "delivery/internal/core/application/usecases/commands/create_order"
	"delivery/internal/generated/queues/basketeventspb"
	"delivery/internal/pkg/errs"

	"github.com/IBM/sarama"
	"google.golang.org/protobuf/proto"
)

type BasketConfirmedConsumer struct {
	consumerGroup sarama.ConsumerGroup
	topic     string
	createOrderHandler createOrder.CreateOrderHandler
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewBasketConfirmedConsumer(brokers []string, groupID string, topic string, createOrderHandler createOrder.CreateOrderHandler) (*BasketConfirmedConsumer, error) {
	if len(brokers) == 0 {
		return nil, errs.NewValueIsRequired("brokers")
	}
	if groupID == "" {
		return nil, errs.NewValueIsRequired("consumerGroup")
	}
	if topic == "" {
		return nil, errs.NewValueIsRequired("topic")
	}
	if createOrderHandler == nil {
		return nil, errs.NewValueIsRequired("createOrderHandler")
	}

	saramaConfig := sarama.NewConfig()
	saramaConfig.Version = sarama.V3_4_0_0
	saramaConfig.Consumer.Return.Errors = true
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest

	kafkaGroup, err := sarama.NewConsumerGroup(brokers, groupID, saramaConfig)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &BasketConfirmedConsumer{
		createOrderHandler: createOrderHandler,
		consumerGroup:      kafkaGroup,
		topic:              topic,
		ctx:                  ctx,
		cancel:               cancel,
	}, nil
}

func (c *BasketConfirmedConsumer) Close() error {
	c.cancel()
	return c.consumerGroup.Close()
}

// Run блокируется и читает топик, пока ctx не отменён или не будет ошибки.
func (c *BasketConfirmedConsumer) Run(ctx context.Context) error {
	for {
		if err := c.consumerGroup.Consume(ctx, []string{c.topic}, c); err != nil {
			return err
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

func (c *BasketConfirmedConsumer) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *BasketConfirmedConsumer) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *BasketConfirmedConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		event := &basketeventspb.BasketConfirmedIntegrationEvent{}
		err := proto.Unmarshal(message.Value, event)
		if err != nil {
			// TODO: deadletters queue?
			session.MarkMessage(message, "")
			continue
		}

		command, err := createOrder.NewCreateOrderCommand(
			event.GetBasketId(),
			event.GetAddress(),
			event.GetItems(),
			event.GetDeliveryPeriod(),
			int(event.GetVolume()),
		)
		if err != nil {
			return err
		}

		if _, err := c.createOrderHandler.Handle(c.ctx, command); err != nil {
			return err
		}

		session.MarkMessage(message, "")
	}

	return nil
}
