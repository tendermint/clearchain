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
	@rm -rf vendor/
	@echo "--> Running glide install"
	@glide install


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
.PHONY: build check_tools get_vendor_deps test benchmark clean
