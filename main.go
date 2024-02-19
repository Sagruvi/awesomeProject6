package main

import (
	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	r.Use(LoggerMiddleware)
	// Применение middleware для логирования с помощью zap logger

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		zap.S().Info(r.RequestURI)
	})
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		zap.S().Info(r.RequestURI)
	})
	r.Put("/", func(w http.ResponseWriter, r *http.Request) {
		zap.S().Info(r.RequestURI)
	})
	http.ListenAndServe(":8080", r)
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger, _ := zap.NewProduction()
		defer logger.Sync()
		logger.Info(r.URL.Path,
			zap.String("RemoteAddr", r.RemoteAddr),
			zap.String("Method", r.Method))
	})
}
