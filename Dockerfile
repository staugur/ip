ARG buildos=golang:1.23-alpine
ARG runos=scratch

# -- build dependencies with alpine --
FROM $buildos AS builder
WORKDIR /build
COPY . .
ARG goproxy
ARG TARGETARCH
RUN if [ "x$goproxy" != "x" ]; then go env -w GOPROXY=${goproxy},direct; fi ; \
    CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -ldflags "-s -w" .

# -- run application with a small image --
FROM $runos
COPY --from=builder /build/mip /bin/
COPY data/ip2region.xdb /
EXPOSE 7000
ENTRYPOINT ["mip", "-db", "/ip2region.xdb"]
