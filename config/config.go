package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

const cfgPath = "./config"

type Config struct {
	AppConfig AppConfig      `mapstructure:"app"`
	APIServer APIServer      `mapstructure:"api_server"`
	Postgres  PostgresConfig `mapstructure:"postgres"`
	Redis     Redis          `mapstructure:"redis"`
	Logger    Logger         `mapstructure:"logger"`
}

type AppConfig struct {
	Name        string `mapstructure:"name"`
	Environment string `mapstructure:"environment"`
	Port        string `mapstructure:"port"`
	Debug       bool   `mapstructure:"debug"`
}

type APIServer struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type PostgresConfig struct {
	Username string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Port     int    `mapstructure:"port"`
	Host     string `mapstructure:"host"`
	DBname   string `mapstructure:"dbname"`
}

type Redis struct {
	Addr     string `mapstructure:"addres"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type Logger struct {
	Level string `mapstructure:"level"`
}

func LoadConfig() (*Config, error) {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AddConfigPath(cfgPath)
	v.SetConfigType("yaml")
	v.SetConfigName("config")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var cfg Config

	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshalling config file: %w", err)
	}

	return &cfg, nil
}

func (c *PostgresConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.Username, c.Password, c.DBname)
}
