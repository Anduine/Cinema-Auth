package app

import (
	"auth-service/internal/config"
	"auth-service/internal/delivery/rest"
	"auth-service/internal/repository"
	"auth-service/internal/server"
	"auth-service/internal/service"
	"auth-service/pkg/database"
	"auth-service/pkg/jwt"
)

func Run(cfg *config.Config) {
	db := database.NewPostgresConnection(&cfg.DB)
	defer db.Close()

	// minioClient := s3db.NewMinioClient(&cfg.S3)

	// rabbitmqClient := rabbitmq.NewRabbitMQClient(&cfg.RabbitMQ)
	// defer rabbitmqClient.Close()

	jwtManager := jwt.NewManager(cfg.Tokens.SecretJWTKey, cfg.Tokens.AccessTokenTTL)

	repo := repository.NewPostgresAuthRepo(db)
	authService := service.NewAuthService(repo, jwtManager)
	handler := rest.NewAuthHandler(authService, cfg.Tokens.AccessTokenTTL, cfg.Tokens.RefreshTokenTTL)

	router := server.NewRouter(handler, jwtManager)

	server.StartServer(cfg.Port, cfg.Timeout, router)
}
