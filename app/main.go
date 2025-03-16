package main

import (
	"context"
	"errors"
	. "labyrinth/logger" // наличие точки перед импортом означает, что я объядиняю  простнаство имен и мне не нужно каждый раз писать logger.Loger, а сразу Loger
	"labyrinth/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Инициализация логера
	log.Println("Initializing logger")
	NewLoger()
	defer Loger.Sync()

	// Инициализация сервера
	log.Println("Initializing server")
	srv := server.NewServer()

	// Канал для graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Запуск сервера
	log.Println("Starting server: http://127.0.0.1:8000")
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("[ ERROR ]: Failed to start server: %v", err)
		}
	}()

	// Ожидание сигнала для graceful shutdown
	<-done
	log.Println("Server is shutting down...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("[ ERROR ]: Failed to shutdown server gracefully: %v", err)
	}

	log.Println("Server has stopped.")
}
