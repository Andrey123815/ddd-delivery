package cmd

import kafkain "delivery/internal/adapters/in/kafka"

func (cr *CompositionRoot) NewBasketConfirmedConsumer() (*kafkain.BasketConfirmedConsumer, error) {
	return kafkain.NewBasketConfirmedConsumer(
		cr.configs.KafkaBrokers(),
		cr.configs.KafkaConsumerGroup,
		cr.configs.KafkaBasketEventsTopic,
		cr.NewCreateOrderHandler(),
	)
}
