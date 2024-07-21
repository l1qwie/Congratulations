FROM golang:1.22.5-bullseye as cert-installer

WORKDIR /app

COPY server.crt /usr/local/share/ca-certificates/server.crt
RUN apt-get update && apt-get install -y ca-certificates && update-ca-certificates

FROM golang:1.22.5-bullseye as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

COPY --from=cert-installer /usr/local/share/ca-certificates/server.crt /usr/local/share/ca-certificates/server.crt
RUN update-ca-certificates

RUN --mount=type=cache,target="/root/.cache/go-build" go build -o bin .

FROM builder as downloader

COPY wait-for-it.sh /usr/local/bin/wait-for-it.sh
RUN chmod +x /usr/local/bin/wait-for-it.sh

FROM downloader as final

CMD ["/app/bin"]
