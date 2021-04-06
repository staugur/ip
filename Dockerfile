# -- build dependencies with alpine --
FROM golang:alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    ip_host=0.0.0.0

WORKDIR /build

COPY . .

RUN go env -w GOPROXY=https://goproxy.cn,direct && \
    go build -ldflags "-s -w" -o ip .

# run application with a small image
FROM scratch

COPY --from=builder /build/ip /bin/

COPY data/ip2region.db /

EXPOSE 7000

ENTRYPOINT ["ip", "-db", "/ip2region.db"]
