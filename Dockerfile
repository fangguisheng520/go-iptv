FROM golang:1.24-alpine AS builder

VOLUME /config
WORKDIR /app
EXPOSE 80 8080

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN  go build -o iptv main.go
RUN chmod +x /app/iptv

FROM alpine:latest
WORKDIR /app

COPY ./client /client
COPY ./apktool/* /usr/bin/
COPY ./static /app/static
COPY ./database /app/database
COPY ./config.yml /app/config.yml
COPY ./README.md  /app/README.md
COPY ./ChangeLog.md /app/ChangeLog.md
COPY ./logo /app/logo
COPY ./Version /app/Version
COPY ./license/license-amd64 /app/license

ENV TZ=Asia/Shanghai
RUN apk add --no-cache openjdk8 bash curl tzdata sqlite;\
    cp /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone; \
    # mkdir -p /config/images/icon ; \
    # mkdir -p /config/images/bj ; \
    chmod 777 -R /usr/bin/apktool* 

COPY --from=builder /app/iptv .

CMD ["./iptv"]