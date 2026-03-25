FROM golang:latest AS builder

ENV GO111MODULE=on\
    GOPROXY=https://goproxy.cn,direct

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o seckill-app main.go

FROM alpine:3.22.3

RUN apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

WORKDIR /app

COPY --from=builder /app/seckill-app .

COPY --from=builder /app/configs ./configs
COPY --from=builder /app/scripts ./scripts

EXPOSE 8080

CMD ["./seckill-app"]