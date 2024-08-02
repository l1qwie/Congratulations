FROM golang:1.22.5-bullseye AS cert-installer

WORKDIR /app

COPY /keys/server.crt /usr/local/share/ca-certificates/server.crt
RUN apt-get update && apt-get install -y ca-certificates && update-ca-certificates

FROM golang:1.22.5-bullseye AS builder

WORKDIR /app

# Сначала скопируем только go.mod и go.sum
COPY go.mod go.sum ./

# Скопируем все модули
COPY Authorization ./Authorization
COPY Employees ./Employees
COPY Notifications ./Notifications
COPY Subscribe ./Subscribe

# Установим зависимости
RUN go mod download

# Скопируем все исходники
COPY . .

# Установим сертификаты
COPY --from=cert-installer /usr/local/share/ca-certificates/server.crt /usr/local/share/ca-certificates/server.crt
RUN update-ca-certificates

# Соберем приложение
RUN --mount=type=cache,target="/root/.cache/go-build" go build -o bin .

FROM builder AS final

CMD ["/app/bin"]
