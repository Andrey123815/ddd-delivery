package kafkaout

import (
	"context"

	"delivery/internal/generated/queues/ordereventspb"
	"delivery/internal/pkg/errs"

	"github.com/IBM/sarama"
	"google.golang.org/protobuf/encoding/protojson"
)

type OrderAssignedProducer struct {
	producer sarama.SyncProducer
	topic    string
}

func NewOrderAssignedProducer(brokers []string, topic string) (*OrderAssignedProducer, error) {
	if len(brokers) == 0 {
		return nil, errs.NewValueIsRequired("brokers")
	}
	if topic == "" {
		return nil, errs.NewValueIsRequired("topic")
	}

	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true // required for SyncProducer (sarama)
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Retry.Max = 10
	cfg.Version = sarama.V3_4_0_0

	producer, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		return nil, err
	}

	return &OrderAssignedProducer{producer: producer, topic: topic}, nil
}

func (p *OrderAssignedProducer) Publish(ctx context.Context, event *ordereventspb.OrderAssignedIntegrationEvent) error {
	if event == nil {
		return errs.NewValueIsRequired("event")
	}
	payload, err := protojson.Marshal(event)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(event.GetOrderId()),
		Value: sarama.ByteEncoder(payload),
		Headers: []sarama.RecordHeader{
			{Key: []byte("content-type"), Value: []byte("application/json")},
		},
	}

	if err := ctx.Err(); err != nil {
		return err
	}
	_, _, err = p.producer.SendMessage(msg)
	return err
}

func (p *OrderAssignedProducer) Close() error {
	return p.producer.Close()
}
