FROM golang:1.21-alpine as builder

ENV GOPROXY=https://goproxy.cn,direct
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# Download dependencies
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# Build
COPY . .
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o app . \
    && rm -rf /app/go.mod /app/go.sum /usr/local/go

FROM alpine as prod

RUN apk add --no-cache docker

WORKDIR /srv/
COPY --from=builder /app/app .
CMD ["./app"]
