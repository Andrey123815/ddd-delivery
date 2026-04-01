package cmd

import "strings"

type Config struct {
	HttpPort               string
	DbHost                 string
	DbPort                 string
	DbUser                 string
	DbPassword             string
	DbName                 string
	DbSslMode              string
	GeoServiceGrpcHost     string
	KafkaHost              string
	KafkaConsumerGroup     string
	KafkaBasketEventsTopic string
	KafkaOrderEventsTopic  string
}

func (c Config) KafkaBrokers() []string {
	parts := strings.Split(c.KafkaHost, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if s := strings.TrimSpace(p); s != "" {
			out = append(out, s)
		}
	}
	return out
}
