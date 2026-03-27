package courierRepo

import (
	"context"
	"delivery/internal/core/domain/models/order"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	uow ports.UnitOfWork
}

func NewRepository(uow ports.UnitOfWork) (*Repository, error) {
	if uow == nil {
		return nil, errs.NewValueIsRequired("uow")
	}

	return &Repository{
		uow: uow,
	}, nil
}

func (r *Repository) Add(ctx context.Context, aggregate *order.Order) error {
	r.uow.Track(aggregate)

  dto := DomainToDTO(aggregate)

	err := r.uow.Tx().WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Create(&dto).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Update(ctx context.Context, aggregate *order.Order) error {
	r.uow.Track(aggregate)

	dto := DomainToDTO(aggregate)

	err := r.uow.Tx().WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Save(&dto).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) Get(ctx context.Context, Id uuid.UUID) (*order.Order, error) {
	dto := OrderDTO{}

	err := r.uow.Tx().WithContext(ctx).Preload(clause.Associations).Find(&dto, Id).Error
	if err != nil {
		return nil, err
	}

	return DtoToDomain(dto), nil
}

func (r *Repository) GetFirstInCreatedStatus(ctx context.Context) (*order.Order, error) {
	dtos := []*OrderDTO{}

	err := r.uow.Tx().WithContext(ctx).
		Preload(clause.Associations).
		Where("status = ?", order.OrderStatusCreated).
		Limit(1).
		Find(&dtos).Error
	if err != nil {
		return nil, err
	}
	if len(dtos) == 0 {
		return nil, nil
	}

	return DtoToDomain(*dtos[0]), nil
}

func (r *Repository) GetAllInAssignedStatus(ctx context.Context) ([]*order.Order, error) {
	dtos := []*OrderDTO{}

	err := r.uow.Tx().WithContext(ctx).Preload(clause.Associations).Where("status = ?", order.OrderStatusAssigned).Find(&dtos).Error
	if err != nil {
		return nil, err
	}

	orders := make([]*order.Order, 0)
	for _, dto := range dtos {
		orders = append(orders, DtoToDomain(*dto))
	}

	return orders, nil
}

func (r *Repository) GetNotCompleted(ctx context.Context) ([]*order.Order, error) {
	dtos := []*OrderDTO{}

	err := r.uow.Tx().WithContext(ctx).
		Preload(clause.Associations).
		Where("status IN ?", []order.OrderStatus{order.OrderStatusCreated, order.OrderStatusAssigned}).
		Find(&dtos).Error
	if err != nil {
		return nil, err
	}

	orders := make([]*order.Order, 0, len(dtos))
	for _, dto := range dtos {
		orders = append(orders, DtoToDomain(*dto))
	}

	return orders, nil
}