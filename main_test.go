package main

import (
	"bytes"
	"fmt"
	"go.uber.org/zap/zapcore"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestLoggerMiddleware(t *testing.T) {
	// Создаем буфер для записи логов
	var buf bytes.Buffer

	// Создаем логгер для тестирования
	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(&buf),
		zap.InfoLevel,
	))
	zap.ReplaceGlobals(logger)
	defer logger.Sync()

	// Создаем фейковый запрос
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Создаем Recorder для записи ответа
	recorder := httptest.NewRecorder()

	// Вызываем middleware с фейковым хэндлером
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	LoggerMiddleware(handler).ServeHTTP(recorder, req)

	// Добавляем задержку для ожидания записи в лог
	time.Sleep(100 * time.Millisecond)

	// Получаем вывод из логгера
	got := buf.String()
	fmt.Printf("got: %s\n", got)

	// Проверяем, что в логе есть запись о запросе
	expectedLog := `{"level":"info","msg":"/","RemoteAddr":"","Method":"GET"}`
	if !strings.Contains(got, expectedLog) {
		t.Errorf("Expected log output %s, got %s", expectedLog, got)
	}
}
