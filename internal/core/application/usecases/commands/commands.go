package commands

import (
	assignCourier "delivery/internal/core/application/usecases/commands/assign_order_to_courier"
	createCourier "delivery/internal/core/application/usecases/commands/create_courier"
	createOrder "delivery/internal/core/application/usecases/commands/create_order"
	moveCouriers "delivery/internal/core/application/usecases/commands/move_couriers"
)

type CreateCourierHandler = createCourier.CreateCourierHandler
type CreateOrderHandler = createOrder.CreateOrderHandler
type AssignCourierHandler = assignCourier.AssignCourierHandler
type MoveCouriersHandler = moveCouriers.MoveCouriersHandler
