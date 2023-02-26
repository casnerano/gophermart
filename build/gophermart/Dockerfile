FROM golang:1.19-alpine3.17 AS builder

RUN apk add --no-cache git

WORKDIR /build

COPY go.mod ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./gophermart ./cmd/gophermart/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /build/gophermart /app
COPY --from=builder /build/configs /app/configs

CMD ["./gophermart"]
