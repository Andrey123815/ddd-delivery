package main

import (
	"context"
	"delivery/cmd"
	httpAdapter "delivery/internal/adapters/in/http"
	"delivery/internal/adapters/in/http/problems"
	"delivery/internal/generated/servers"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/robfig/cron/v3"
)

func main() {
	config := getConfigs()

	compositionRoot := cmd.NewCompositionRoot(config)
	defer compositionRoot.CloseAll()

	startCron(compositionRoot)
	startKafkaConsumer(compositionRoot)
	startWebServer(compositionRoot, config.HttpPort)
}

func getConfigs() cmd.Config {
	return cmd.Config{
		HttpPort:               goDotEnvVariable("HTTP_PORT"),
		DbHost:                 goDotEnvVariable("DB_HOST"),
		DbPort:                 goDotEnvVariable("DB_PORT"),
		DbUser:                 goDotEnvVariable("DB_USER"),
		DbPassword:             goDotEnvVariable("DB_PASSWORD"),
		DbName:                 goDotEnvVariable("DB_NAME"),
		DbSslMode:              goDotEnvVariable("DB_SSLMODE"),
		GeoServiceGrpcHost:     goDotEnvVariable("GEO_SERVICE_GRPC_HOST"),
		KafkaHost:              goDotEnvVariable("KAFKA_HOST"),
		KafkaConsumerGroup:     goDotEnvVariable("KAFKA_CONSUMER_GROUP"),
		KafkaBasketEventsTopic: goDotEnvVariable("KAFKA_BASKET_EVENTS_TOPIC"),
		KafkaOrderEventsTopic:  goDotEnvVariable("KAFKA_ORDER_EVENTS_TOPIC"),
	}
}

func goDotEnvVariable(key string) string {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

func startCron(cr *cmd.CompositionRoot) {
	c := cron.New()

	if _, err := c.AddJob("@every 1s", cr.NewAssignOrdersJob()); err != nil {
		log.Fatalf("ошибка при добавлении AssignOrdersJob: %v", err)
	}

	if _, err := c.AddJob("@every 1s", cr.NewMoveCouriersJob()); err != nil {
		log.Fatalf("ошибка при добавлении MoveCouriersJob: %v", err)
	}

	c.Start()
}

func startWebServer(cr *cmd.CompositionRoot, port string) {
	httpServer, err := httpAdapter.NewServer(
		cr.NewCreateCourierHandler(),
		cr.NewCreateOrderHandler(),
		cr.NewGetAllCouriersHandler(),
		cr.NewGetNoCompletedOrdersHandler(),
	)
	if err != nil {
		log.Fatalf("cannot create http server: %v", err)
	}

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
	}))
	e.HTTPErrorHandler = problems.EchoErrorHandler

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Healthy")
	})

	servers.RegisterHandlers(e, httpServer)

	e.Logger.Fatal(e.Start(fmt.Sprintf("0.0.0.0:%s", port)))
}

func startKafkaConsumer(cr *cmd.CompositionRoot) {
	consumer, err := cr.NewBasketConfirmedConsumer()
	if err != nil {
		log.Fatalf("cannot create BasketConfirmedConsumer: %v", err)
	}
	cr.RegisterCloser(consumer)

	go func() {
		if err := consumer.Run(context.Background()); err != nil {
			log.Printf("kafka consumer stopped: %v", err)
		}
	}()
}
