FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o gateway main.go

FROM alpine:3.24
RUN apk --no-cache add ca-certificates curl
WORKDIR /root/
COPY --from=builder /app/gateway .

EXPOSE 8080
CMD ["./gateway"]
