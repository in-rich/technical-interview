COVER_FILE=$(CURDIR)/coverage.out
BIN_DIR=$(CURDIR)/bin

PKG="github.com/in-rich/technical-interview/server"

PKG_LIST=$(shell go list $(PKG)/... | grep -v /vendor/)

# Runs the test suite.
test:
	go run gotest.tools/gotestsum@latest --packages="./..." --junitfile report.xml --format pkgname -- -count=1 -p 1 -v -coverpkg=./...

firebase:
	firebase emulators:start

# Starts the database emulator.
run-db:
	direnv allow . && source .envrc && firebase emulators:start

# Starts the development server.
run:
	direnv allow . && source .envrc && go run ./cmd/server/main.go
