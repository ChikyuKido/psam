FROM golang:1.24-alpine as builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .
ENV CGO_ENABLED=1
RUN go build -o psam .

FROM alpine:latest

WORKDIR /app/

COPY --from=builder /app/psam .

EXPOSE 8080

CMD ["./psam","serve"]
