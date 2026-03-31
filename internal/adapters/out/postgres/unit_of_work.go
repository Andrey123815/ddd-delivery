package postgres

import (
	"context"
	courierRepo "delivery/internal/adapters/out/postgres/courierRepo"
	outboxRepo "delivery/internal/adapters/out/postgres/outboxRepo"
	orderRepo "delivery/internal/adapters/out/postgres/orderRepo"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/ddd"
	"delivery/internal/pkg/errs"
	"errors"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

var _ ports.UnitOfWork = (*UnitOfWork)(nil)

type UnitOfWork struct {
	tx                *gorm.DB
	db                *gorm.DB
	committed         bool
	trackedAggregates []ddd.AggregateRoot
	courierRepository ports.CourierRepository
	orderRepository   ports.OrderRepository
	outboxRepository  ports.OutboxRepository
}

func NewUnitOfWork(db *gorm.DB) (ports.UnitOfWork, error) {
	if db == nil {
		return nil, errs.NewValueIsRequired("db")
	}

	uow := &UnitOfWork{
		db: db,
	}

	courierRepo, err := courierRepo.NewRepository(uow)
	if err != nil {
		return nil, err
	}
	uow.courierRepository = courierRepo

	orderRepo, err := orderRepo.NewRepository(uow)
	if err != nil {
		return nil, err
	}
	uow.orderRepository = orderRepo

	outRepo, err := outboxRepo.NewRepository(uow)
	if err != nil {
		return nil, err
	}
	uow.outboxRepository = outRepo

	return uow, nil
}

func (u *UnitOfWork) Tx() *gorm.DB {
	return u.tx
}

func (u *UnitOfWork) Db() *gorm.DB {
	return u.db
}

func (u *UnitOfWork) InTx() bool {
	return u.tx != nil
}

func (u *UnitOfWork) Track(agg ddd.AggregateRoot) {
	u.trackedAggregates = append(u.trackedAggregates, agg)
}

func (u *UnitOfWork) CourierRepository() ports.CourierRepository {
	return u.courierRepository
}

func (u *UnitOfWork) OrderRepository() ports.OrderRepository {
	return u.orderRepository
}

func (u *UnitOfWork) OutboxRepository() ports.OutboxRepository {
	return u.outboxRepository
}

func (u *UnitOfWork) Begin(ctx context.Context) {
	u.tx = u.db.WithContext(ctx).Begin()
	u.committed = false
}

func (u *UnitOfWork) Commit(ctx context.Context) error {
	if u.tx == nil {
		return errs.NewValueIsRequired("cannot commit without transaction")
	}

	if err := u.tx.WithContext(ctx).Commit().Error; err != nil {
		return err
	}

	u.committed = true
	u.clearTx()
	return nil
}

func (u *UnitOfWork) RollbackUnlessCommitted(ctx context.Context) {
	if u.tx != nil && !u.committed {
		if err := u.tx.WithContext(ctx).Rollback().Error; err != nil && !errors.Is(err, gorm.ErrInvalidTransaction) {
			log.Error(err)
		}
		u.clearTx()
	}
}

func (u *UnitOfWork) clearTx() {
	u.tx = nil
	u.trackedAggregates = nil
	u.committed = false
}
