package utils

import (
	"delivery/cmd"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetConfigs() cmd.Config {
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

