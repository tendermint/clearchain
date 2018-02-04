PACKAGES=$(shell go list ./... | grep -v '/vendor/')
BUILD_FLAGS = -ldflags "-X github.com/tendermint/clearchain/version.GitCommit=`git rev-parse --short HEAD`"

all: get_vendor_deps build test

########################################
### Build

build:
	go build $(BUILD_FLAGS) -o build/clearchain ./cmd/clearchain


########################################
### Tools & dependencies

get_vendor_deps:
	@rm -rf vendor/
	@echo "--> Running glide install"
	@glide install


########################################
### Testing

test:
	@go test -v -cover $(PACKAGES)

benchmark:
	@go test -bench=. $(PACKAGES)


# To avoid unintended conflicts with file names, always add to .PHONY
# unless there is a reason not to.
# https://www.gnu.org/software/make/manual/html_node/Phony-Targets.html
.PHONY: build check_tools get_vendor_deps test benchmark
