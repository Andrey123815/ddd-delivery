package cmd

import kafkaout "delivery/internal/adapters/out/kafka"

func (cr *CompositionRoot) NewOrderAssignedProducer() (*kafkaout.OrderAssignedProducer, error) {
	return kafkaout.NewOrderAssignedProducer(
		cr.configs.KafkaBrokers(),
		cr.configs.KafkaOrderEventsTopic,
	)
}

func (cr *CompositionRoot) NewOrderCompletedProducer() (*kafkaout.OrderCompletedProducer, error) {
	return kafkaout.NewOrderCompletedProducer(
		cr.configs.KafkaBrokers(),
		cr.configs.KafkaOrderEventsTopic,
	)
}
