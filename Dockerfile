ARG buildos=golang:1.17.0-alpine
ARG runos=scratch

# -- build dependencies with alpine --
FROM $buildos AS builder

WORKDIR /build

COPY . .

ARG TARGETARCH

RUN go env -w GOPROXY=https://goproxy.cn,direct && \
    CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -ldflags "-s -w" .

# run application with a small image
FROM $runos

COPY --from=builder /build/ip /bin/

COPY data/ip2region.db /

EXPOSE 7000

ENTRYPOINT ["ip", "-db", "/ip2region.db"]
