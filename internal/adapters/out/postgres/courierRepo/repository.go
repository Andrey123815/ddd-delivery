package courierRepo

import (
	"context"
	"delivery/internal/core/domain/models/courier"
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

func (r *Repository) Add(ctx context.Context, aggregate *courier.Courier) error {
	r.uow.Track(aggregate)

  dto := DomainToDTO(aggregate)

	err := r.uow.Tx().WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Create(&dto).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Update(ctx context.Context, aggregate *courier.Courier) error {
	r.uow.Track(aggregate)

	dto := DomainToDTO(aggregate)

	err := r.uow.Tx().WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Save(&dto).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) Get(ctx context.Context, Id uuid.UUID) (*courier.Courier, error) {
	dto := CourierDTO{}

	err := r.uow.Tx().WithContext(ctx).Preload(clause.Associations).Find(&dto, Id).Error
	if err != nil {
		return nil, err
	}

	return DtoToDomain(dto), nil
}

func (r *Repository) GetAllAvailableCouriers(ctx context.Context) ([]*courier.Courier, error) {
	dtos := []*CourierDTO{}

	err := r.uow.Tx().WithContext(ctx).Preload(clause.Associations).Find(&dtos).Error
	if err != nil {
		return nil, err
	}

	couriers := make([]*courier.Courier, 0)
	for _, dto := range dtos {
		courier := DtoToDomain(*dto)

		courierIsAvailable := true

		for _, storagePlace := range courier.StoragePlaces() {
			if storagePlace.IsOccupied() {
				courierIsAvailable = false
				break
			}
		}

		if courierIsAvailable == false {
			continue
		}

		couriers = append(couriers, DtoToDomain(*dto))
	}

	return couriers, nil
}
