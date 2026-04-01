package cmd

import (
	"context"
	"delivery/internal/adapters/out/postgres"
	"delivery/internal/jobs"
	"log"

	"github.com/robfig/cron/v3"
)

///////////////////////////////////////////////////////////
//////////////////// JOBS /////////////////////////////////
///////////////////////////////////////////////////////////

func (cr *CompositionRoot) NewAssignOrdersJob() cron.Job	 {
	job, err := jobs.NewAssignOrdersJob(cr.NewAssignCourierHandler())
	if err != nil {
		log.Fatalf("cannot create AssignOrdersJob: %v", err)
	}

	return job
}

func (cr *CompositionRoot) NewMoveCouriersJob() cron.Job {
	job, err := jobs.NewMoveCouriersJob(cr.NewMoveCouriersHandler())
	if err != nil {
		log.Fatalf("cannot create MoveCouriersJob: %v", err)
	}

	return job
}

func (cr *CompositionRoot) NewOutboxProcessorJob() cron.Job {
	ctx := context.Background()
	
	factory, err := postgres.NewUnitOfWorkFactory(cr.DB())
	if err != nil {
		log.Fatalf("cannot create UnitOfWorkFactory: %v", err)
	}

	uow, err := factory.New(ctx)
	if err != nil {
		log.Fatalf("cannot create UnitOfWork: %v", err)
	}

	job, err := jobs.NewProcessOutboxMessagesJob(uow.OutboxRepository(), cr.EventRegistry)
	if err != nil {
		log.Fatalf("cannot create ProcessOutboxMessagesJob: %v", err)
	}

	return job
}