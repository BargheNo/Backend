FROM golang:1.23.4-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

RUN go env -w GOPROXY=https://goproxy.io,direct

COPY go.mod go.sum ./

RUN go mod download

COPY . .
COPY ./internal/infrastructure/jwt/privateKey.pem ./internal/infrastructure/jwt/
COPY ./internal/infrastructure/jwt/publicKey.pem /internal/infrastructure/jwt/

COPY ./internal/infrastructure/communication/email/templates ./internal/infrastructure/communication/email/templates


RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o main ./cmd/app

FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/main .

RUN mkdir -p /app/internal/infrastructure/jwt
RUN mkdir -p /app/internal/infrastructure/communication/email/templates
COPY ./internal/infrastructure/jwt/privateKey.pem ./internal/infrastructure/jwt/
COPY ./internal/infrastructure/jwt/publicKey.pem ./internal/infrastructure/jwt/
COPY ./internal/infrastructure/communication/email/templates ./internal/infrastructure/communication/email/templates

COPY .env .


EXPOSE 8080

CMD ["./main"]