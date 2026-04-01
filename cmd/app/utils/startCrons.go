package utils

import (
	"delivery/cmd"
	"log"

	"github.com/robfig/cron/v3"
)

func StartCrons(cr *cmd.CompositionRoot) {
	c := cron.New()

	if _, err := c.AddJob("@every 1s", cr.NewAssignOrdersJob()); err != nil {
		log.Fatalf("ошибка при добавлении AssignOrdersJob: %v", err)
	}

	if _, err := c.AddJob("@every 1s", cr.NewMoveCouriersJob()); err != nil {
		log.Fatalf("ошибка при добавлении MoveCouriersJob: %v", err)
	}

	if _, err := c.AddJob("@every 1s", cr.NewOutboxProcessorJob()); err != nil {
		log.Fatalf("cannot add OutboxProcessorCronJob: %v", err)
	}

	c.Start()
}
