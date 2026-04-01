package cmd

import (
	"fmt"
	"log"
	"sync"

	"delivery/internal/adapters/out/outbox"
	courierRepo "delivery/internal/adapters/out/postgres/courierRepo"
	orderRepo "delivery/internal/adapters/out/postgres/orderRepo"
	postgresgorm "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbOnce     sync.Once
	dbInstance *gorm.DB
)

func (cr *CompositionRoot) DB() *gorm.DB {
	dbOnce.Do(func() {
		dsn := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cr.configs.DbHost,
			cr.configs.DbPort,
			cr.configs.DbUser,
			cr.configs.DbPassword,
			cr.configs.DbName,
			cr.configs.DbSslMode,
		)

		db, err := gorm.Open(postgresgorm.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("cannot connect to database: %v", err)
		}

		if err := db.AutoMigrate(
			&courierRepo.CourierDTO{},
			&courierRepo.StoragePlaceDTO{},
			&orderRepo.OrderDTO{},
			&outbox.Message{},
		); err != nil {
			log.Fatalf("cannot run auto migrate: %v", err)
		}

		dbInstance = db
	})

	return dbInstance
}
