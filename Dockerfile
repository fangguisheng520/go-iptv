FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o iptv .


FROM eclipse-temurin:8-jre-alpine
WORKDIR /app

COPY --from=builder /app/iptv .

RUN chmod +x ./iptv

CMD ["./iptv","-port=8080","-conf=/config","-build=/build","-java=/opt/java/openjdk/bin"]