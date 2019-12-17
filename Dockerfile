# Используем базовый образ для Go
FROM golang:latest

# Создадим директорию
RUN mkdir /app

# Скопируем всё в директорию
COPY . .

# Получим зависимости, которые использовали в боте
RUN go get github.com/sergejkoll/tg-botkp2019
RUN go get github.com/Syfaro/telegram-bot-api

# Соберём приложение
RUN go build *.go

# Запустим приложение
CMD go run *.go
