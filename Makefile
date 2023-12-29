APP_NAME=dot

default: build

all: clean test build

clean:
	@ echo "+ $@"
	rm -rf ./bin

# -timeout 	timout in seconds
#  -v		verbose output
test:
	@ echo "+ $@"
	go test -timeout=5s -v

# `CGO_ENABLED=0`
# Because of dynamically linked libraries, this will statically compile the
# app with all libraries built in. You won't be able to cross-compile if CGO
# is enabled. This is because Go binary is looking for libraries on the
# operating system it’s running in. We compiled our app, but it still is
# dynamically linked to the libraries it needs to run
# (i.e., all the C libraries it binds to). When using a minimal docker image
# the operating system doesn't have these libraries.
#
# `-a`
# Force rebuilding of package, all import will be rebuilt with cgo disabled,
# which means all the imports will be rebuilt with cgo disabled.
#
# `-installsuffix cgo`
# A suffix to use in the name of the package installation directory
#
# `-o`
# Output
#
# `./bin/[name-of-app]`
# Placement of the binary
#
# `.`
# Location of the source files
build:
	@ echo "+ $@"
	CGO_ENABLED=0 go build -a -installsuffix cgo -o ./bin/${APP_NAME} .

# Cross-compile
# http://dave.cheney.net/2015/08/22/cross-compilation-with-go-1-5
build-linux:
	@ echo "+ $@"
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -o ./bin/${APP_NAME}-linux-amd64 .

build-mac:
	@ echo "+ $@"
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -o ./bin/${APP_NAME}-darwin-amd64 .

.PHONY: all default clean test build build-linux build-mac
