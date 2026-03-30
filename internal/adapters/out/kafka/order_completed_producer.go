package kafkaout

import (
	"context"
	"log"

	"delivery/internal/generated/queues/ordereventspb"
	"delivery/internal/pkg/errs"

	"github.com/IBM/sarama"
	"google.golang.org/protobuf/encoding/protojson"
)

type OrderCompletedProducer struct {
	producer sarama.AsyncProducer
	topic    string
}

func NewOrderCompletedProducer(brokers []string, topic string) (*OrderCompletedProducer, error) {
	if len(brokers) == 0 {
		return nil, errs.NewValueIsRequired("brokers")
	}
	if topic == "" {
		return nil, errs.NewValueIsRequired("topic")
	}

	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Retry.Max = 10
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true
	cfg.Version = sarama.V3_4_0_0

	producer, err := sarama.NewAsyncProducer(brokers, cfg)
	if err != nil {
		return nil, err
	}

	go func() {
		for range producer.Successes() {
		}
	}()
	go func() {
		for pe := range producer.Errors() {
			log.Printf("OrderCompletedProducer: kafka error: %v", pe.Err)
		}
	}()

	return &OrderCompletedProducer{producer: producer, topic: topic}, nil
}

func (p *OrderCompletedProducer) Publish(ctx context.Context, event *ordereventspb.OrderCompletedIntegrationEvent) error {
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

	select {
	case p.producer.Input() <- msg:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (p *OrderCompletedProducer) Close() error {
	return p.producer.Close()
}
