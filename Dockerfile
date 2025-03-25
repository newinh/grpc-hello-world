FROM golang:1.23 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download -x

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server/main.go

FROM golang:1.23-alpine
WORKDIR /app

COPY --from=builder /app/main ./main

ENTRYPOINT ["./main"]
