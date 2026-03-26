package uow

import (
	"delivery/internal/core/ports"
	"delivery/internal/pkg/ddd"

	"gorm.io/gorm"
)

type UnitOfWork interface {
	Tx() *gorm.DB
	Db() *gorm.DB
	InTx() bool
	Begin()
	Commit() error
	Track(agg ddd.AggregateRoot)
	CourierRepository() ports.CourierRepository
	OrderRepository() ports.OrderRepository
}
