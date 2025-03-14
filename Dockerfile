FROM golang:1.23.4-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

RUN go env -w GOPROXY=https://goproxy.io,direct

COPY go.mod go.sum ./

RUN go mod download

COPY . .
COPY ./internal/application/adapter/jwt/privateKey.pem ./internal/application/adapter/jwt/
COPY ./internal/application/adapter/jwt/publicKey.pem /internal/application/adapter/jwt/

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main ./cmd/app

FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/main .

RUN mkdir -p /app/internal/application/adapter/jwt
COPY ./internal/application/adapter/jwt/privateKey.pem ./internal/application/adapter/jwt/
COPY ./internal/application/adapter/jwt/publicKey.pem ./internal/application/adapter/jwt/
COPY .env .

EXPOSE 8080

CMD ["./main"]