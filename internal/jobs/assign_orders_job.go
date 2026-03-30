package jobs

import (
	"context"
	"delivery/internal/core/application/usecases/commands"
	"delivery/internal/pkg/errs"

	"github.com/labstack/gommon/log"
	"github.com/robfig/cron/v3"
)

var _ cron.Job = &AssignOrdersJob{}

type AssignOrdersJob struct {
	handler commands.AssignCourierHandler
}

func NewAssignOrdersJob(handler commands.AssignCourierHandler) (cron.Job, error) {
	if handler == nil {
		return nil, errs.NewValueIsRequired("handler")
	}
	return &AssignOrdersJob{handler: handler}, nil
}

func (j *AssignOrdersJob) Run() {
	if err := j.handler.Handle(context.Background()); err != nil {
		log.Error(err)
	}
}
