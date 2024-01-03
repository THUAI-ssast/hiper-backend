FROM golang:1.21-alpine as builder

ENV GOPROXY=https://goproxy.cn,direct

# Download dependencies
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# Build
COPY . .
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o app . \
    && rm -rf /app/go.mod /app/go.sum /usr/local/go

FROM alpine as prod

apk add --no-cache docker

WORKDIR /srv/
COPY --from=builder /app/app .
CMD ["./app"]
