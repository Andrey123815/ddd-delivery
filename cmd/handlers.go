package cmd

import (
	"delivery/internal/adapters/out/grpc/geo"
	"delivery/internal/adapters/out/postgres"
	"delivery/internal/core/application/usecases/commands"
	assignCourier "delivery/internal/core/application/usecases/commands/assign_order_to_courier"
	createCourier "delivery/internal/core/application/usecases/commands/create_courier"
	createOrder "delivery/internal/core/application/usecases/commands/create_order"
	moveCouriers "delivery/internal/core/application/usecases/commands/move_couriers"
	"delivery/internal/core/application/usecases/queries"
	getAllCouriers "delivery/internal/core/application/usecases/queries/get_all_couriers"
	getNoCompletedOrders "delivery/internal/core/application/usecases/queries/get_no_completed_orders"
	"delivery/internal/core/domain/services"
	"log"
)

func (cr *CompositionRoot) NewCreateCourierHandler() commands.CreateCourierHandler {
	factory, err := postgres.NewUnitOfWorkFactory(cr.DB())
	if err != nil {
		log.Fatalf("cannot create UnitOfWorkFactory: %v", err)
	}

	handler, err := createCourier.NewCreateCourierHandler(factory)
	if err != nil {
		log.Fatalf("cannot create CreateCourierHandler: %v", err)
	}
	
	return handler
}

func (cr *CompositionRoot) NewCreateOrderHandler() commands.CreateOrderHandler {
	factory, err := postgres.NewUnitOfWorkFactory(cr.DB())
	if err != nil {
		log.Fatalf("cannot create UnitOfWorkFactory: %v", err)
	}

	geoClient, err := geo.NewClient(cr.configs.GeoServiceGrpcHost)
	if err != nil {
		log.Fatalf("cannot create GeoClient: %v", err)
	}

	handler, err := createOrder.NewCreateOrderHandler(factory, geoClient)
	if err != nil {
		log.Fatalf("cannot create CreateOrderHandler: %v", err)
	}
	
	return handler
}

func (cr *CompositionRoot) NewAssignCourierHandler() commands.AssignCourierHandler {
	factory, err := postgres.NewUnitOfWorkFactory(cr.DB())
	if err != nil {
		log.Fatalf("cannot create UnitOfWorkFactory: %v", err)
	}

	dispatcher := services.NewOrderDispatcher()

	handler, err := assignCourier.NewAssignCourierHandler(
		dispatcher,
		factory,
	)
	if err != nil {
		log.Fatalf("cannot create AssignCourierHandler: %v", err)
	}

	return handler
}

func (cr *CompositionRoot) NewMoveCouriersHandler() commands.MoveCouriersHandler {
	factory, err := postgres.NewUnitOfWorkFactory(cr.DB())
	if err != nil {
		log.Fatalf("cannot create UnitOfWorkFactory: %v", err)
	}

	handler, err := moveCouriers.NewMoveCouriersHandler(factory)
	if err != nil {
		log.Fatalf("cannot create MoveCouriersHandler: %v", err)
	}

	return handler
}

func (cr *CompositionRoot) NewGetAllCouriersHandler() queries.GetAllCouriersHandler {
	factory, err := postgres.NewUnitOfWorkFactory(cr.DB())
	if err != nil {
		log.Fatalf("cannot create UnitOfWorkFactory: %v", err)
	}

	handler, err := getAllCouriers.NewGetAllCouriersHandler(factory)
	if err != nil {
		log.Fatalf("cannot create GetAllCouriersHandler: %v", err)
	}

	return handler
}

func (cr *CompositionRoot) NewGetNoCompletedOrdersHandler() queries.GetNoCompletedOrdersHandler {
	factory, err := postgres.NewUnitOfWorkFactory(cr.DB())
	if err != nil {
		log.Fatalf("cannot create UnitOfWorkFactory: %v", err)
	}

	handler, err := getNoCompletedOrders.NewGetNoCompletedOrdersHandler(factory)
	if err != nil {
		log.Fatalf("cannot create GetNoCompletedOrdersHandler: %v", err)
	}

	return handler
}
