package server

import (
	"net/http"
	"time"
)

func NewServer() *http.Server {
	return &http.Server{
		Handler:      NewRouter(),
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}
