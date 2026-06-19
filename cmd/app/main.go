package main

import (
	"auth-service/internal/app"
	"auth-service/internal/config"
	"auth-service/pkg/logger"
)

func main() {
	config := config.MustLoadConfig()

	logger.SetupGlobalLogger(config.Logger.Env, config.Logger.LogLevel, config.Logger.LogTimeZone)

	app.Run(config)
}
