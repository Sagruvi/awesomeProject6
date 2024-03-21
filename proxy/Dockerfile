# Используем официальный образ Go как базовый
FROM golang:1.18-alpine

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app


# Копируем исходники приложения в рабочую директорию
COPY . .
COPY ./docs . /usr/local/go/src/main/proxy/
COPY ./graph . /usr/local/go/src/main/proxy/
COPY ./internal . /usr/local/go/src/main/proxy/


# Скачиваем все зависимости
RUN go mod init  test && go mod tidy

# Собираем приложение
RUN go build -o main ./cmd

# Открываем порт 8080
EXPOSE 8080

# Запускаем приложение
CMD ["/main"]
