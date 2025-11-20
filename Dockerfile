FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# бинарник API
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/api/main.go

# бинарник мигратора
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o migrate ./cmd/migrate/main.go

FROM alpine:3.18

WORKDIR /app
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/app .
COPY --from=builder /app/migrate .
COPY migrations ./migrations

EXPOSE 8080

CMD ["./app"]
