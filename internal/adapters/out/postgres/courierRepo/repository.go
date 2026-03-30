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

	return r.uow.Tx().WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Create(&dto).Error
}

func (r *Repository) Update(ctx context.Context, aggregate *courier.Courier) error {
	r.uow.Track(aggregate)

	dto := DomainToDTO(aggregate)
	tx := r.uow.Tx().WithContext(ctx)

	return tx.Session(&gorm.Session{FullSaveAssociations: true}).Save(&dto).Error
}

func (r *Repository) Get(ctx context.Context, Id uuid.UUID) (*courier.Courier, error) {
	dto := CourierDTO{}

	err := r.uow.Tx().WithContext(ctx).Preload(clause.Associations).Find(&dto, Id).Error
	if err != nil {
		return nil, err
	}

	return DtoToDomain(dto), nil
}

func (r *Repository) GetAll(ctx context.Context) ([]*courier.Courier, error) {
	dtos := []*CourierDTO{}
	if err := r.uow.Tx().WithContext(ctx).Preload(clause.Associations).Find(&dtos).Error; err != nil {
		return nil, err
	}

	out := make([]*courier.Courier, 0, len(dtos))
	for _, dto := range dtos {
		out = append(out, DtoToDomain(*dto))
	}
	
	return out, nil
}

func (r *Repository) GetAllAvailableCouriers(ctx context.Context) ([]*courier.Courier, error) {
	var dtos []*CourierDTO
	if err := r.uow.Tx().WithContext(ctx).Preload(clause.Associations).Find(&dtos).Error; err != nil {
		return nil, err
	}

	freeCouriers := make([]*courier.Courier, 0, len(dtos))
	for _, dto := range dtos {
		c := DtoToDomain(*dto)
		if c.IsFree() {
			freeCouriers = append(freeCouriers, c)
		}
	}
	return freeCouriers, nil
}
