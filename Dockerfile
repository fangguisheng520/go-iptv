FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN  go build -o iptv main.go


FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/iptv .
COPY ./client /client
COPY ./apktool/* /usr/bin/
COPY ./static /app/static
COPY ./database /app/database
COPY ./config.yml /app/config.yml
COPY ./README.md  /app/README.md
COPY ./logo /app/logo

ENV TZ=Asia/Shanghai
RUN apk add --no-cache openjdk8 bash curl tzdata sqlite;\
    cp /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone; \
    mkdir -p /config/images/icon ; \
    mkdir -p /config/images/bj ; \
    chmod 777 -R /app/iptv /usr/bin/apktool* 

CMD ["./iptv"]