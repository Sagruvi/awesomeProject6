# Используем официальный образ Go как базовый
FROM golang:1.18-alpine as builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app


# Копируем исходники приложения в рабочую директорию
COPY . .
COPY ./geolocation /usr/local/go/src/main/proxy/geolocation
COPY ./docs /usr/local/go/src/main/proxy/docs
COPY ./internal /usr/local/go/src/main/proxy/internal

# Скачиваем все зависимости
RUN go mod init  test && go mod tidy

# Собираем приложение
RUN go build -o main ./cmd

# Начинаем новую стадию сборки на основе минимального образа
FROM alpine:latest

# Добавляем исполняемый файл из первой стадии в корневую директорию контейнера
COPY --from=builder /app/main /main

# Открываем порт 8080
EXPOSE 8080

# Запускаем приложение
CMD ["/main"]
