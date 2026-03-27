package jobs

import (
	"context"
	"delivery/internal/core/application/usecases/commands"
	"delivery/internal/pkg/errs"

	"github.com/labstack/gommon/log"
	"github.com/robfig/cron/v3"
)

var _ cron.Job = &MoveCouriersJob{}

type MoveCouriersJob struct {
	handler commands.MoveCouriersHandler
}

func NewMoveCouriersJob(handler commands.MoveCouriersHandler) (cron.Job, error) {
	if handler == nil {
		return nil, errs.NewValueIsRequired("handler")
	}
	return &MoveCouriersJob{handler: handler}, nil
}

func (j *MoveCouriersJob) Run() {
	if err := j.handler.Handle(context.Background()); err != nil {
		log.Errorf("MoveCouriersJob: %v", err)
	}
}
