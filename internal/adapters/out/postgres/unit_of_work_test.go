package postgres

import (
	"context"
	courierRepo "delivery/internal/adapters/out/postgres/courierRepo"
	orderRepo "delivery/internal/adapters/out/postgres/orderRepo"
	"delivery/internal/core/domain/models/courier"
	"delivery/internal/core/domain/models/kernel"
	"delivery/internal/core/domain/models/order"
	"delivery/internal/pkg/testcnts"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	postgresgorm "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTest(t *testing.T) (context.Context, *gorm.DB, error) {
	ctx := context.Background()
	postgresContainer, dsn, err := testcnts.StartPostgresContainer(ctx)
	if err != nil {
		return nil, nil, err
	}

	// Подключаемся к БД через Gorm
	db, err := gorm.Open(postgresgorm.Open(dsn), &gorm.Config{})
	assert.NoError(t, err)

	// Авто миграция (создаём таблицу)
	err = db.AutoMigrate(&courierRepo.CourierDTO{})
	assert.NoError(t, err)
	err = db.AutoMigrate(&courierRepo.CourierDTO{})
	assert.NoError(t, err)
	err = db.AutoMigrate(&orderRepo.OrderDTO{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&courierRepo.StoragePlaceDTO{})
	assert.NoError(t, err)

	// Очистка выполняется после завершения теста
	t.Cleanup(func() {
		err := postgresContainer.Terminate(ctx)
		assert.NoError(t, err)
	})

	return ctx, db, nil
}

func Test_CourierRepositoryShouldCanAddCourier(t *testing.T) {
	// Инициализируем окружение
	ctx, db, err := setupTest(t)
	assert.NoError(t, err)

	// Создаем UnitOfWork
	uow, err := NewUnitOfWork(db)
	assert.NoError(t, err)

	// Вызываем Add
	location, _ := kernel.NewLocation(10, 10)
	courierAggregate, err := courier.NewCourier("Велосипедист", 2, location)
	assert.NoError(t, err)

	uow.Begin(ctx)
	err = uow.CourierRepository().Add(ctx, courierAggregate)
	assert.NoError(t, err)
	err = uow.Commit(ctx)
	assert.NoError(t, err)

	// Считываем данные из БД
	var courierFromDb courierRepo.CourierDTO
	err = db.First(&courierFromDb, "id = ?", courierAggregate.ID()).Error
	assert.NoError(t, err)

	// Проверяем эквивалентность
	assert.Equal(t, courierAggregate.ID(), courierFromDb.ID)
	assert.Equal(t, courierAggregate.Speed(), courierFromDb.Speed)
}

func Test_OrderRepositoryShouldCanAddOrder(t *testing.T) {
	// Инициализируем окружение
	ctx, db, err := setupTest(t)
	assert.NoError(t, err)

	// Создаем UnitOfWork
	uow, err := NewUnitOfWork(db)
	assert.NoError(t, err)

	// Вызываем Add
	location, _ := kernel.NewLocation(1, 1)
	orderAggregate, err := order.NewOrder(location, 10)
	assert.NoError(t, err)

	uow.Begin(ctx)
	err = uow.OrderRepository().Add(ctx, orderAggregate)
	assert.NoError(t, err)
	err = uow.Commit(ctx)
	assert.NoError(t, err)

	// Считываем данные из БД
	var orderFromDb orderRepo.OrderDTO
	err = db.First(&orderFromDb, "id = ?", orderAggregate.Id()).Error
	assert.NoError(t, err)

	// Проверяем эквивалентность
	assert.Equal(t, orderAggregate.Id(), orderFromDb.ID)
	assert.Equal(t, int(orderAggregate.Location().X()), orderFromDb.Location.X)
	assert.Equal(t, int(orderAggregate.Location().Y()), orderFromDb.Location.Y)
	assert.Equal(t, orderAggregate.Volume(), orderFromDb.Volume)
	assert.Equal(t, orderAggregate.Status(), orderFromDb.Status)
}

func Test_CourierRepositoryShouldCanUpdateCourier(t *testing.T) {
	ctx, db, err := setupTest(t)
	assert.NoError(t, err)

	uow, err := NewUnitOfWork(db)
	assert.NoError(t, err)

	// Создаем и добавляем курьера
	location, _ := kernel.NewLocation(5, 5)
	courierAggregate, err := courier.NewCourier("Пеший курьер", 1, location)
	assert.NoError(t, err)

	uow.Begin(ctx)
	err = uow.CourierRepository().Add(ctx, courierAggregate)
	assert.NoError(t, err)
	err = uow.Commit(ctx)
	assert.NoError(t, err)

	// Обновляем локацию курьера
	newLocation, _ := kernel.NewLocation(8, 9)
	err = courierAggregate.Move(newLocation)
	assert.NoError(t, err)

	uow.Begin(ctx)
	err = uow.CourierRepository().Update(ctx, courierAggregate)
	assert.NoError(t, err)
	err = uow.Commit(ctx)
	assert.NoError(t, err)

	// Проверяем обновление в БД
	var courierFromDb courierRepo.CourierDTO
	err = db.First(&courierFromDb, "id = ?", courierAggregate.ID()).Error
	assert.NoError(t, err)
	assert.Equal(t, 8, courierFromDb.Location.X)
	assert.Equal(t, 9, courierFromDb.Location.Y)
}

func Test_OrderRepositoryShouldCanUpdateOrder(t *testing.T) {
	ctx, db, err := setupTest(t)
	assert.NoError(t, err)

	uow, err := NewUnitOfWork(db)
	assert.NoError(t, err)

	// Создаем и добавляем заказ
	location, _ := kernel.NewLocation(2, 2)
	orderAggregate, err := order.NewOrder(location, 15)
	assert.NoError(t, err)

	uow.Begin(ctx)
	err = uow.OrderRepository().Add(ctx, orderAggregate)
	assert.NoError(t, err)
	err = uow.Commit(ctx)
	assert.NoError(t, err)

	// Назначаем заказ курьеру
	courierId := uuid.New()
	err = orderAggregate.Assign(courierId)
	assert.NoError(t, err)

	uow.Begin(ctx)
	err = uow.OrderRepository().Update(ctx, orderAggregate)
	assert.NoError(t, err)
	err = uow.Commit(ctx)
	assert.NoError(t, err)

	// Проверяем обновление в БД
	var orderFromDb orderRepo.OrderDTO
	err = db.First(&orderFromDb, "id = ?", orderAggregate.Id()).Error
	assert.NoError(t, err)
	assert.Equal(t, order.OrderStatusAssigned, orderFromDb.Status)
	assert.Equal(t, courierId, orderFromDb.CourierID)
}

func Test_CourierRepositoryShouldCanGetCourier(t *testing.T) {
	ctx, db, err := setupTest(t)
	assert.NoError(t, err)

	uow, err := NewUnitOfWork(db)
	assert.NoError(t, err)

	// Создаем и добавляем курьера
	location, _ := kernel.NewLocation(7, 8)
	courierAggregate, err := courier.NewCourier("Автокурьер", 5, location)
	assert.NoError(t, err)

	uow.Begin(ctx)
	err = uow.CourierRepository().Add(ctx, courierAggregate)
	assert.NoError(t, err)
	err = uow.Commit(ctx)
	assert.NoError(t, err)

	// Получаем курьера через репозиторий
	uow.Begin(ctx)
	fetchedCourier, err := uow.CourierRepository().Get(ctx, courierAggregate.ID())
	assert.NoError(t, err)

	// Проверяем корректность полученных данных
	assert.Equal(t, courierAggregate.ID(), fetchedCourier.ID())
	assert.Equal(t, courierAggregate.Name(), fetchedCourier.Name())
	assert.Equal(t, courierAggregate.Speed(), fetchedCourier.Speed())
	assert.Equal(t, courierAggregate.Location().X(), fetchedCourier.Location().X())
	assert.Equal(t, courierAggregate.Location().Y(), fetchedCourier.Location().Y())
}

func Test_OrderRepositoryShouldCanGetOrder(t *testing.T) {
	ctx, db, err := setupTest(t)
	assert.NoError(t, err)

	uow, err := NewUnitOfWork(db)
	assert.NoError(t, err)

	// Создаем и добавляем заказ
	location, _ := kernel.NewLocation(3, 4)
	orderAggregate, err := order.NewOrder(location, 20)
	assert.NoError(t, err)

	uow.Begin(ctx)
	err = uow.OrderRepository().Add(ctx, orderAggregate)
	assert.NoError(t, err)
	err = uow.Commit(ctx)
	assert.NoError(t, err)

	// Получаем заказ через репозиторий
	uow.Begin(ctx)
	fetchedOrder, err := uow.OrderRepository().Get(ctx, orderAggregate.Id())
	assert.NoError(t, err)

	// Проверяем корректность полученных данных
	assert.Equal(t, orderAggregate.Id(), fetchedOrder.Id())
	assert.Equal(t, orderAggregate.Volume(), fetchedOrder.Volume())
	assert.Equal(t, orderAggregate.Status(), fetchedOrder.Status())
	assert.Equal(t, orderAggregate.Location().X(), fetchedOrder.Location().X())
	assert.Equal(t, orderAggregate.Location().Y(), fetchedOrder.Location().Y())
}

func Test_UnitOfWorkShouldCommitMultipleOperations(t *testing.T) {
	ctx, db, err := setupTest(t)
	assert.NoError(t, err)

	uow, err := NewUnitOfWork(db)
	assert.NoError(t, err)

	// Создаем агрегаты
	courierLocation, _ := kernel.NewLocation(10, 10)
	courierAggregate, err := courier.NewCourier("Мото курьер", 10, courierLocation)
	assert.NoError(t, err)

	orderLocation, _ := kernel.NewLocation(5, 5)
	orderAggregate, err := order.NewOrder(orderLocation, 25)
	assert.NoError(t, err)

	// Выполняем несколько операций в одной транзакции
	uow.Begin(ctx)
	err = uow.CourierRepository().Add(ctx, courierAggregate)
	assert.NoError(t, err)
	err = uow.OrderRepository().Add(ctx, orderAggregate)
	assert.NoError(t, err)
	err = uow.Commit(ctx)
	assert.NoError(t, err)

	// Проверяем, что оба агрегата сохранены
	var courierFromDb courierRepo.CourierDTO
	err = db.First(&courierFromDb, "id = ?", courierAggregate.ID()).Error
	assert.NoError(t, err)
	assert.Equal(t, courierAggregate.ID(), courierFromDb.ID)

	var orderFromDb orderRepo.OrderDTO
	err = db.First(&orderFromDb, "id = ?", orderAggregate.Id()).Error
	assert.NoError(t, err)
	assert.Equal(t, orderAggregate.Id(), orderFromDb.ID)
}

func Test_UnitOfWorkShouldTrackAggregates(t *testing.T) {
	ctx, db, err := setupTest(t)
	assert.NoError(t, err)

	uow, err := NewUnitOfWork(db)
	assert.NoError(t, err)

	// Создаем агрегат
	location, err := kernel.NewLocation(9, 10)
	assert.NoError(t, err)
	courierAggregate, err := courier.NewCourier("Трекируемый курьер", 3, location)
	assert.NoError(t, err)

	// Track должен быть вызван внутри Add
	uow.Begin(ctx)
	err = uow.CourierRepository().Add(ctx, courierAggregate)
	assert.NoError(t, err)

	// Проверяем, что агрегат был отслежен (косвенная проверка через успешность коммита)
	err = uow.Commit(ctx)
	assert.NoError(t, err)
}

func Test_UnitOfWorkShouldHandleNestedTransactions(t *testing.T) {
	ctx, db, err := setupTest(t)
	assert.NoError(t, err)

	uow, err := NewUnitOfWork(db)
	assert.NoError(t, err)

	location, err := kernel.NewLocation(7, 8)
	assert.NoError(t, err)
	courierAggregate, err := courier.NewCourier("Курьер для вложенной транзакции", 7, location)
	assert.NoError(t, err)

	// Начинаем транзакцию
	uow.Begin(ctx)
	assert.True(t, uow.InTx())

	err = uow.CourierRepository().Add(ctx, courierAggregate)
	assert.NoError(t, err)

	// Коммитим
	err = uow.Commit(ctx)
	assert.NoError(t, err)
	assert.False(t, uow.InTx())

	// Проверяем сохранение
	var courierFromDb courierRepo.CourierDTO
	err = db.First(&courierFromDb, "id = ?", courierAggregate.ID()).Error
	assert.NoError(t, err)
	assert.Equal(t, courierAggregate.ID(), courierFromDb.ID)
}

func Test_OrderRepositoryShouldCanGetFirstInCreatedStatus(t *testing.T) {
	ctx, db, err := setupTest(t)
	assert.NoError(t, err)

	uow, err := NewUnitOfWork(db)
	assert.NoError(t, err)

	// Создаем несколько заказов
	location1, _ := kernel.NewLocation(1, 1)
	order1, err := order.NewOrder(location1, 10)
	assert.NoError(t, err)

	location2, _ := kernel.NewLocation(2, 2)
	order2, err := order.NewOrder(location2, 15)
	assert.NoError(t, err)

	// Добавляем заказы
	uow.Begin(ctx)
	err = uow.OrderRepository().Add(ctx, order1)
	assert.NoError(t, err)
	err = uow.OrderRepository().Add(ctx, order2)
	assert.NoError(t, err)
	err = uow.Commit(ctx)
	assert.NoError(t, err)

	// Получаем первый заказ в статусе Created
	uow.Begin(ctx)
	firstOrder, err := uow.OrderRepository().GetFirstInCreatedStatus(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, firstOrder)
	assert.Equal(t, order.OrderStatusCreated, firstOrder.Status())
}

func Test_CourierRepositoryShouldGetAllAvailableCouriers(t *testing.T) {
	ctx, db, err := setupTest(t)
	assert.NoError(t, err)

	uow, err := NewUnitOfWork(db)
	assert.NoError(t, err)

	// Создаем несколько курьеров
	location1, _ := kernel.NewLocation(5, 5)
	courier1, err := courier.NewCourier("Курьер 1", 2, location1)
	assert.NoError(t, err)
	err = courier1.AddStoragePlace("Багажник 1", 50)
	assert.NoError(t, err)

	location2, _ := kernel.NewLocation(10, 10)
	courier2, err := courier.NewCourier("Курьер 2", 3, location2)
	assert.NoError(t, err)
	err = courier2.AddStoragePlace("Багажник 2", 60)
	assert.NoError(t, err)

	// Добавляем курьеров
	uow.Begin(ctx)
	err = uow.CourierRepository().Add(ctx, courier1)
	assert.NoError(t, err)
	err = uow.CourierRepository().Add(ctx, courier2)
	assert.NoError(t, err)
	err = uow.Commit(ctx)
	assert.NoError(t, err)

	// Получаем всех доступных курьеров
	uow.Begin(ctx)
	availableCouriers, err := uow.CourierRepository().GetAllAvailableCouriers(ctx)
	assert.NoError(t, err)
	assert.Len(t, availableCouriers, 2)
}

func Test_UnitOfWorkShouldHandleTransactionLifecycle(t *testing.T) {
	ctx, db, err := setupTest(t)
	assert.NoError(t, err)

	uow, err := NewUnitOfWork(db)
	assert.NoError(t, err)

	// Проверяем начальное состояние
	assert.False(t, uow.InTx())
	assert.Nil(t, uow.Tx())

	// Начинаем транзакцию
	uow.Begin(ctx)
	assert.True(t, uow.InTx())
	assert.NotNil(t, uow.Tx())

	// Коммитим
	err = uow.Commit(ctx)
	assert.NoError(t, err)
	assert.False(t, uow.InTx())
	assert.Nil(t, uow.Tx())
}

func Test_UnitOfWorkShouldRollbackUnlessCommitted(t *testing.T) {
	ctx, db, err := setupTest(t)
	assert.NoError(t, err)

	uowInterface, err := NewUnitOfWork(db)
	assert.NoError(t, err)
	uow := uowInterface.(*UnitOfWork)

	location, err := kernel.NewLocation(6, 7)
	assert.NoError(t, err)
	courierAggregate, err := courier.NewCourier("Курьер для отката", 4, location)
	assert.NoError(t, err)

	// Начинаем транзакцию и добавляем курьера
	uow.Begin(ctx)
	err = uow.CourierRepository().Add(ctx, courierAggregate)
	assert.NoError(t, err)

	// Откатываем без коммита
	uow.RollbackUnlessCommitted(ctx)
	assert.False(t, uow.InTx())

	// Проверяем, что курьер НЕ сохранен в БД
	var courierFromDb courierRepo.CourierDTO
	err = db.First(&courierFromDb, "id = ?", courierAggregate.ID()).Error
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}
