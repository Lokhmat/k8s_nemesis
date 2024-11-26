package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/transport"
	"github.com/go-kit/kit/transport/http"
)

// Service определяет интерфейс сервиса
type Service interface {
	Hello(ctx context.Context) (string, error)
}

// HelloService реализует интерфейс Service
type HelloService struct{}

// Hello возвращает приветственное сообщение
func (HelloService) Hello(ctx context.Context) (string, error) {
	return "Hello, World!", nil
}

// MakeHelloEndpoint создает эндпоинт для метода Hello
func MakeHelloEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		msg, err := svc.Hello(ctx)
		return map[string]string{"message": msg}, err
	}
}

// decodeRequest декодирует HTTP-запрос
func decodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return struct{}{}, nil
}

// encodeResponse кодирует ответ в JSON
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

func main() {
	// Настройка логгера
	logger := log.NewLogfmtLogger(os.Stdout)
	logger = level.NewFilter(logger, level.AllowAll())
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	// Инициализация сервиса
	var svc Service
	svc = HelloService{}

	// Создание эндпоинта
	helloEndpoint := MakeHelloEndpoint(svc)

	// Создание HTTP-обработчика
	handler := http.NewServer(
		helloEndpoint,
		decodeRequest,
		encodeResponse,
		transport.ServerErrorLogger(logger),
	)

	// Регистрация маршрута
	http.Handle("/hello", handler)

	// Запуск HTTP-сервера
	level.Info(logger).Log("msg", "Starting server on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		level.Error(logger).Log("msg", "Server failed", "error", err)
	}
}
