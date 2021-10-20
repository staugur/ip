.PHONY: help

BINARY=mip
Version=$(shell grep "const version" main.go | tr -d '"' | awk '{print $$NF}')

help:
	@echo "  make clean  - Remove binaries and vim swap files"
	@echo "  make gotool - Run go tool 'fmt' and 'vet'"
	@echo "  make build  - Compile go code and generate binary file"
	@echo "  make dev    - Run dev server"
	@echo "  make release- Format go code and compile to generate binary release"

gotool:
	go fmt ./
	go vet ./

build: gotool
	go build -ldflags "-s -w" -o bin/$(BINARY) && chmod +x bin/$(BINARY)

docker:
	docker build -t staugur/ip .

dev:
	@echo Starting service...
	@go run ./

release: build
	mv bin/$(BINARY) . && tar zcvf bin/$(BINARY).$(Version)-linux-amd64.tar.gz $(BINARY) data && rm $(BINARY)
