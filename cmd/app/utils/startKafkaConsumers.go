package utils

import (
	"context"
	"delivery/cmd"
	"log"
)

func StartKafkaConsumers(cr *cmd.CompositionRoot) {
	consumer, err := cr.NewBasketConfirmedConsumer()
	if err != nil {
		log.Fatalf("cannot create BasketConfirmedConsumer: %v", err)
	}
	cr.RegisterCloser(consumer)

	go func() {
		if err := consumer.Run(context.Background()); err != nil {
			log.Printf("kafka consumer stopped: %v", err)
		}
	}()
}
