default: test-unit

# build the binary
build:
	go install

# Fetch base dependencies as well as testing packages
get-deps:
	go get github.com/golang/lint/golint
	# Ginkgo and omega test tools
	go get github.com/onsi/ginkgo/ginkgo
	go get github.com/onsi/gomega

# Cleans up directory and source code with gofmt
clean:
	go clean ./...

# Run gofmt on all code
fmt:
	gofmt -l -w $$(ls -d */ | grep -v vendor)

# Run linter with non-strict checking
lint:
	ls -d */ | grep -v vendor | xargs -L 1 golint

# Vet code
vet:
	go tool vet $$(ls -d */ | grep -v vendor)

# Perform only unit tests
test-unit: get-deps clean fmt lint vet build
	ginkgo -r -skipPackage integration

help:
	 @echo "common developer commands:"
	 @echo "  get-deps: fetch developer dependencies"
	 @echo "  fmt: run gofmt on the codebase"
	 @echo "  clean: run go clean on the codebase"
	 @echo "  lint: run go lint on the codebase"
	 @echo "  vet: run go vet on the codebase"
	 @echo ""
	 @echo "common testing commands:"
	 @echo "  test-unit: Unit tests"
