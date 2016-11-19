mkfile_path := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

export GOPATH := $(mkfile_path)/go
export PATH := $(mkfile_path)/bin:$(PATH)
export TMROOT = $(mkfile_path)/.tendermint/


build: $(GOPATH)/bin/tendermint

apps:
	go get -u github.com/tendermint/tmsp/cmd/...

init: build
	tendermint init

$(GOPATH)/bin/tendermint: glide
	go get github.com/tendermint/tendermint/cmd/tendermint
	cd $(GOPATH)/src/github.com/tendermint/tendermint ; \
	glide install ; go install ./cmd/tendermint

glide:
	go get github.com/Masterminds/glide

.PHONY: build init glide apps
