# Используем официальный образ Go как базовый
FROM golang:1.19-alpine

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем исходники приложения в рабочую директорию
COPY . .

# Скачиваем все зависимости
RUN go mod init main && \
    go mod tidy
RUN mkdir proxy
COPY . ./proxy
# Собираем приложение
RUN go build -o main ./proxy/cmd

# Открываем порт 8080
EXPOSE 8080

# Запускаем приложение
CMD ["./main"]
