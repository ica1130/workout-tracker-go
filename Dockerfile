FROM golang:1.22 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o workout-tracker ./cmd/api

FROM alpine:3.18

ENV ENV_PATH=/app/.env

WORKDIR /app

COPY --from=builder /app/workout-tracker .
COPY .env .env

EXPOSE 4000

CMD ["./workout-tracker"]
