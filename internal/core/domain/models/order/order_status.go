package order

type OrderStatus string

const (
	OrderStatusCreated  OrderStatus = "Created"
	OrderStatusAssigned   OrderStatus = "Assigned"
	OrderStatusCompleted  OrderStatus = "Completed"
)
