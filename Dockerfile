FROM golang:1.23.2
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o bot ./cmd/main.go
CMD ["./bot"]