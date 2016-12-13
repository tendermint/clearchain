NAME = ledger
HARDWARE = $(shell uname -m)
VERSION ?= 0.0.2
BUILD_TAG ?= dev
BUILD_DIR ?= build

CWD := $(shell pwd)
GOLINT := $(GOPATH)/bin/golint
LDFLAGS_DEFAULT = -X=app.Version=$(VERSION)

ci: $(BUILD_DIR)/linux/amd64

all: $(BUILD_DIR)/linux/amd64 $(BUILD_DIR)/linux/386 \
	$(BUILD_DIR)/darwin/amd64 $(BUILD_DIR)/darwin/386 \
	$(BUILD_DIR)/windows/amd64 $(BUILD_DIR)/windows/386

darwin: $(BUILD_DIR)/darwin/amd64 $(BUILD_DIR)/darwin/386
	lipo -create build/darwin/386/ledger build/darwin/amd64/ledger -output build/darwin/ledger

$(BUILD_DIR)/%: deps
	GOOS=$(word 2,$(subst /, ,$@)) GOARCH=$(word 3,$(subst /, ,$@)) \
	go build -v -ldflags="$(LDFLAGS_DEFAULT) $(LDFLAGS)" -o $@/$(NAME) ./cmd/$(NAME)
	GOOS=$(word 2,$(subst /, ,$@)) GOARCH=$(word 3,$(subst /, ,$@)) \
	go build -v -ldflags="$(LDFLAGS_DEFAULT) $(LDFLAGS)" -o $@/ledgerctl ./cmd/ledgerctl
	GOOS=$(word 2,$(subst /, ,$@)) GOARCH=$(word 3,$(subst /, ,$@)) \
	go build -v -ldflags="$(LDFLAGS_DEFAULT) $(LDFLAGS)" -o $@/webserver ./cmd/webserver

deps: $(BUILD_DIR)/deps-stamp
$(BUILD_DIR)/deps-stamp:
	go get -u -v github.com/golang/lint/golint
	go get -u -v github.com/golang/mock/mockgen
	go get -d -v ./... || true
	mkdir -p $(BUILD_DIR)
	touch $@

release:

mockgen:
	mockgen -package mock_tx github.com/tenermint/clearchain/ledger/types Tx,TxExecutor > testutil/mocks/mock_tx/mock_tx.go

test: deps
	go test -race -cover -v ./...

lint:
	go vet ./... || true
	$(GOLINT) ./... || true

dist-clean: clean
	rm -f $(BUILD_DIR)/deps-stamp

clean:
	rm -rf $(BUILD_DIR)/*/ -rf

.PHONY: release deps clean test lint
