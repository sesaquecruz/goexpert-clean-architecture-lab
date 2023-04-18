package config

import "github.com/spf13/viper"

type config struct {
	DBDriver       string `mapstructure:"DB_DRIVER"`
	DBHost         string `mapstructure:"DB_HOST"`
	DBPort         string `mapstructure:"DB_PORT"`
	DBUser         string `mapstructure:"DB_USER"`
	DBPassword     string `mapstructure:"DB_PASSWORD"`
	DBName         string `mapstructure:"DB_NAME"`
	AMQPHost       string `mapstructure:"AMQP_HOST"`
	AMQPPort       string `mapstructure:"AMQP_PORT"`
	AMQPUser       string `mapstructure:"AMQP_USER"`
	AMQPPassword   string `mapstructure:"AMQP_PASSWORD"`
	GRPCServerPort string `mapstructure:"GRPC_SERVER_PORT"`
}

func LoadConfig(path string) (*config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	var cfg *config
	if err := viper.Unmarshal(&cfg); err != nil {
		panic(err)
	}

	return cfg, nil
}
