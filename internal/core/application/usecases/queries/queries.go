package queries

import (
	getAllCouriers "delivery/internal/core/application/usecases/queries/get_all_couriers"
	getNoCompletedOrders "delivery/internal/core/application/usecases/queries/get_no_completed_orders"
)

type GetAllCouriersHandler = getAllCouriers.GetAllCouriersHandler
type GetNoCompletedOrdersHandler = getNoCompletedOrders.GetNoCompletedOrdersHandler
