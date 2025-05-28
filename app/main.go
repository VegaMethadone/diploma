package main

import (
	"context"
	"fmt"
	"labyrinth/logger"
	"labyrinth/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	fmt.Println("Starting server...")

	currentTime := time.Now()
	dateDir := currentTime.Format("02_01_2006")
	timeFile := currentTime.Format("15_04")
	loggerPath := fmt.Sprintf("../logs/%s/%s.log", dateDir, timeFile)
	if err := os.MkdirAll(fmt.Sprintf("../logs/%s", dateDir), 0755); err != nil {
		panic(fmt.Sprintf("Failed to create log directory: %v", err))
	}
	fmt.Println("[ LOGS PATH ]: ", loggerPath)
	logger.InitFileLogger(loggerPath)
	logger.NewInfoMessage("Server is starting...",
		zap.Time("time", currentTime),
	)

	// Настройка HTTP-сервера с таймаутами
	httpServer := server.NewServer()

	// Канал для graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Запуск сервера в горутине
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	fmt.Printf("Server started on %s\n", httpServer.Addr)

	// Ожидание сигнала завершения
	<-done
	fmt.Println("\nServer is shutting down...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	fmt.Println("Server stopped")
}
