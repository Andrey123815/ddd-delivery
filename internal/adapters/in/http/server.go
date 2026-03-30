package http

import (
	"delivery/internal/core/application/usecases/commands"
	"delivery/internal/core/application/usecases/queries"
	"delivery/internal/generated/servers"
	"delivery/internal/pkg/errs"
)

var _ servers.ServerInterface = &Server{}

type Server struct {
	createCourierHandler        commands.CreateCourierHandler
	createOrderHandler          commands.CreateOrderHandler
	getAllCouriersHandler        queries.GetAllCouriersHandler
	getNoCompletedOrdersHandler queries.GetNoCompletedOrdersHandler
}

func NewServer(
	createCourierHandler commands.CreateCourierHandler,
	createOrderHandler commands.CreateOrderHandler,
	getAllCouriersHandler queries.GetAllCouriersHandler,
	getNoCompletedOrdersHandler queries.GetNoCompletedOrdersHandler,
) (*Server, error) {
	if createCourierHandler == nil {
		return nil, errs.NewValueIsRequired("createCourierHandler")
	}
	if createOrderHandler == nil {
		return nil, errs.NewValueIsRequired("createOrderHandler")
	}
	if getAllCouriersHandler == nil {
		return nil, errs.NewValueIsRequired("getAllCouriersHandler")
	}
	if getNoCompletedOrdersHandler == nil {
		return nil, errs.NewValueIsRequired("getNoCompletedOrdersHandler")
	}

	return &Server{
		createCourierHandler:        createCourierHandler,
		createOrderHandler:          createOrderHandler,
		getAllCouriersHandler:        getAllCouriersHandler,
		getNoCompletedOrdersHandler: getNoCompletedOrdersHandler,
	}, nil
}
