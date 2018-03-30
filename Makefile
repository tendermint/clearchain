PACKAGES=$(shell go list ./... | grep -v '/vendor/')
BUILD_FLAGS = -ldflags "-X github.com/tendermint/clearchain.Version=`git describe`"

all: get_vendor_deps build test

########################################
### Build

build: clearchaind

clearchaind:
	go build $(BUILD_FLAGS) ./cmd/clearchaind


########################################
### Tools & dependencies

get_vendor_deps:
	@echo "--> Purge old vendor/ directory and run dep ensure"
	rm -rf vendor/ ; dep ensure -v

dep: $(GOPATH)/bin/dep
$(GOPATH)/bin/dep:
	[ -d $(GOPATH)/bin ] || mkdir $(GOPATH)/bin
	@wget https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 -O $@ && chmod +x $@


########################################
### Testing

test: coverage.txt
coverage.txt: clean
	touch $@
	for p in $(PACKAGES); do \
	  rm -f profile.out ;\
	  go test -v -race -coverprofile=profile.out -covermode=atomic $$p;\
	  [ ! -f profile.out ] || \
	    ( cat profile.out >> $@ ; rm profile.out ) \
	done

clean:
	rm -f profile.out coverage.txt

benchmark:
	@go test -bench=. $(PACKAGES)


# To avoid unintended conflicts with file names, always add to .PHONY
# unless there is a reason not to.
# https://www.gnu.org/software/make/manual/html_node/Phony-Targets.html
.PHONY: build dep get_vendor_deps test benchmark clean
