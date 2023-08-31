FROM golang:latest

WORKDIR /app

COPY init.sql /docker-entrypoint-initdb.d/

RUN go install github.com/cosmtrek/air@latest

COPY go.mod go.sum ./
RUN go mod download

CMD ["air", "-c", ".air.toml"]
