

GO           ?= go
GOFMT        ?= $(GO)fmt
FIRST_GOPATH := $(firstword $(subst :, ,$(shell $(GO) env GOPATH)))

BIN_DIR ?= $(shell pwd)/build
VERSION ?= $(shell cat VERSION)


all:  fmt style  build 
 
style:
	@echo ">> checking code style"
	! $(GOFMT) -d $$(find . -name '*.go' -print) | grep '^'

build: 
	@echo ">> building vsphere-exporter binaries"
	$(GO) build -o build/vsphere-exporter main.go

fmt:
	@echo ">> format code style"
	$(GOFMT) -w $$(find . -name '*.go' -print) 

clean:
	rm -rf $(BIN_DIR)

.PHONY: all style  fmt  build clean