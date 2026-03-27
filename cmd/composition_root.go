package cmd

import (
	"delivery/internal/jobs"
	"log"

	"github.com/robfig/cron/v3"
)

type CompositionRoot struct {
	configs Config

	closers []Closer
}

func NewCompositionRoot(configs Config) *CompositionRoot {
	return &CompositionRoot{
		configs: configs,
	}
}

///////////////////////////////////////////////////////////
//////////////////// JOBS /////////////////////////////////
///////////////////////////////////////////////////////////

func (cr *CompositionRoot) NewAssignOrdersJob() cron.Job {
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

///////////////////////////////////////////////////////////
//////////////////// LIFECYCLE ////////////////////////////
///////////////////////////////////////////////////////////

func (cr *CompositionRoot) RegisterCloser(c Closer) {
	cr.closers = append(cr.closers, c)
}

func (cr *CompositionRoot) CloseAll() {
	for _, c := range cr.closers {
		if err := c.Close(); err != nil {
			log.Printf("close error: %v", err)
		}
	}
}
