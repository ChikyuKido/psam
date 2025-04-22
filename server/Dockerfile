FROM golang:1.24-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN go build -o psam .

FROM alpine:latest

WORKDIR /app/

COPY --from=builder /app/psam .

EXPOSE 8080

CMD ["./psam"]
