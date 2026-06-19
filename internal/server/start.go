package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func StartServer(port string, timeout time.Duration, router http.Handler) {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("Сервіс auth працює на порту: " + port)
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			slog.Error("Помилка сервера: ", slog.Any("error", err))
		}
	}()

	<-sigChan
	slog.Info("Вимкнення сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	slog.Info("Shutdown ", "stopcode", server.Shutdown(ctx))
}
