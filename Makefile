# https://github.com/johannesboyne/godockersample
# https://blog.codeship.com/building-minimal-docker-containers-for-go-applications/
# https://developer.atlassion.com/blog/2015/07/osx-static-golang-binaries-with-docker/
DEPS = \
	github.com/codegangsta/cli \
	github.com/fatih/color

default: test

# -timeout	timeout in seconds
# -v		verbose output
test:
	go test -timeout=5s -v

# -d	download but don't install
# -v	verbose output
deps:
	go get -d -v $(DEPS)

# `CGO_ENABLED=0`
# Because of dynamically linked libraries, this will statically compile the
# app with all libraries built in. You won't be able to cross-compile if CGO
# is enabled.
#
# `GOOS=linux`
#
# `-a`
# Force rebuilding of packages, all import will be rebuilt with cgo disabled.
#
# `-installsuffix cgo`
#  A suffix to use in the name of the package installation directory
#
# `-o`
# Output
#
# `bin/dot`
# Placement of the binary
build: deps
	@mkdir -p bin/
	CGO_ENABLED=1 go build -a -installsuffix cgo -o bin/dot


.PHONY: default test deps build
