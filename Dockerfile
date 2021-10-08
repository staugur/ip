ARG buildos=golang:1.17.2-alpine
ARG runos=scratch

# -- build dependencies with alpine --
FROM $buildos AS builder

LABEL maintainer=me@tcw.im

WORKDIR /build

COPY . .

ARG goproxy

ARG TARGETARCH=amd64

RUN if [ "x$goproxy" != "x" ]; then go env -w GOPROXY=${goproxy},direct; fi ; \
    CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -ldflags "-s -w" .

# -- run application with a small image --
FROM $runos

COPY --from=builder /build/ip /bin/

COPY data/ip2region.db /

EXPOSE 7000

ENTRYPOINT ["ip", "-db", "/ip2region.db"]
