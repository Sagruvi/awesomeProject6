package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {
	// Создаем реверсивный прокси для сервера Hugo
	hugoURL, err := url.Parse("http://hugo:1313")
	if err != nil {
		panic(err)
	}
	hugoProxy := httputil.NewSingleHostReverseProxy(hugoURL)

	// Обработчик для маршрута /api/
	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from API"))
	})

	// Обработчик для всех остальных запросов
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, если запрос начинается с /api/
		if strings.HasPrefix(r.URL.Path, "/api/") {
			// Если да, отправляем текст "Hello from API"
			w.Write([]byte("Hello from API"))
		} else {
			// Если нет, перенаправляем запрос на сервер Hugo
			hugoProxy.ServeHTTP(w, r)
		}
	})
	WorkerTest()
	// Запускаем сервер на порту 8080
	http.ListenAndServe(":8080", nil)
}

type ReverseProxy struct {
	host string
	port string
}

func NewReverseProxy(host, port string) *ReverseProxy {
	return &ReverseProxy{
		host: host,
		port: port,
	}
}

func (rp *ReverseProxy) ReverseProxy(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}

const content = ``

func WorkerTest() {
	t := time.NewTicker(1 * time.Second)
	var b byte = 0
	for {
		select {
		case <-t.C:
			err := writeToFile("./app/static/_index.md", fmt.Sprintf("%s%d", content, b))
			if err != nil {
				log.Fatal(err)
			}
			b++
		}
	}
}
func writeToFile(path string, data string) error {
	// Открываем файл для записи, флаг O_TRUNC обрежет файл до нуля и удалит предыдущее содержимое
	file, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Записываем данные в файл
	_, err = file.WriteString(data)
	if err != nil {
		return err
	}

	return nil
}
