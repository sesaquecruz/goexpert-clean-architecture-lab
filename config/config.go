package config

import (
	"fmt"
	"os"
	"reflect"
)

type Config struct {
	DBDriver          string `env:"DB_DRIVER"`
	DBHost            string `env:"DB_HOST"`
	DBPort            string `env:"DB_PORT"`
	DBUser            string `env:"DB_USER"`
	DBPassword        string `env:"DB_PASSWORD"`
	DBName            string `env:"DB_NAME"`
	AMQPHost          string `env:"AMQP_HOST"`
	AMQPPort          string `env:"AMQP_PORT"`
	AMQPUser          string `env:"AMQP_USER"`
	AMQPPassword      string `env:"AMQP_PASSWORD"`
	GRPCServerPort    string `env:"GRPC_SERVER_PORT"`
	RESTServerPort    string `env:"REST_SERVER_PORT"`
	GRAPHQLServerPort string `env:"GRAPHQL_SERVER_PORT"`
}

func LoadConfig() (*Config, error) {
	config := Config{}

	valueOf := reflect.ValueOf(&config).Elem()
	typeOf := valueOf.Type()

	for i := 0; i < valueOf.NumField(); i++ {
		field := valueOf.Field(i)
		envName := typeOf.Field(i).Tag.Get("env")
		envValue, ok := os.LookupEnv(envName)

		if !ok {
			return nil, fmt.Errorf("%s was not found", envName)
		}

		field.SetString(envValue)
	}

	return &config, nil
}
