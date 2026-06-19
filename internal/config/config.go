package config

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Port          string        `yaml:"port" env:"PORT" env-default:"8080"`
	Timeout       time.Duration `yaml:"timeout" env:"TIMEOUT" env-default:"5s"`
	ApiGatewayURL string        `yaml:"api_gateway_url" env:"API_GATEWAY_URL"`
	Logger        LoggerConfig  `yaml:"logger"`
	Tokens        TokensConfig  `yaml:"tokens"`
	DB            DBConfig      `yaml:"db"`
}

type LoggerConfig struct {
	Env         string `yaml:"env" env:"ENV" env-default:"local"`
	LogLevel    string `yaml:"log_level" env:"LOG_LEVEL" env-default:"info"`
	LogTimeZone int    `yaml:"log_time_zone" env:"LOG_TIME_ZONE" env-default:"UTC"`
}

type TokensConfig struct {
	SecretJWTKey    string        `yaml:"secret_jwt_key" env:"SECRET_JWT_KEY"`
	AccessTokenTTL  time.Duration `yaml:"access_token_ttl" env:"ACCESS_TOKEN_TTL"`
	RefreshTokenTTL time.Duration `yaml:"refresh_token_ttl" env:"REFRESH_TOKEN_TTL"`
}

type DBConfig struct {
	Host     string `yaml:"host" env:"DB_HOST" env-required:"true"`
	DBName   string `yaml:"db_name" env:"DB_NAME" env-required:"true"`
	User     string `yaml:"user" env:"DB_USER" env-required:"true"`
	Password string `yaml:"password" env:"DB_PASSWORD" env-required:"true"`
}

func MustLoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Failed to load env: " + err.Error())
	}

	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist: " + path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("Failed to load config: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var result string

	// --config="path/to/config.yaml"
	flag.StringVar(&result, "config", "", "path to config file")
	flag.Parse()

	if result == "" {
		result = os.Getenv("CONFIG_PATH")
	}

	if result == "" {
		result = "./configs/config.yaml"
	}

	return result
}
