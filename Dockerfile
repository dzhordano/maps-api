FROM golang:1.23.4-alpine3.20 AS builder

RUN mkdir /app
ADD . /app
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/app/main.go

FROM alpine:3.20

RUN mkdir /app

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/.env .

RUN mkdir /docs
COPY --from=builder /app/docs ./docs

EXPOSE 9000

CMD ["./main"]