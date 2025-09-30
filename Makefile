# Build

APP_NAME=gcopy
DESCRIPTION=Cross-platform file copy utility
VERSION=1.0.0
BUILD_TIME=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-X 'main.Version=$(VERSION)' -X 'main.Description=$(DESCRIPTION)' -X 'main.BuildTime=$(BUILD_TIME)'

all: windows linux mac-darwin mac-arm

windows:
	GOOS=windows GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o build/$(APP_NAME)-windows-amd64.exe ./cmd/gcopy/

linux:
	GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o build/$(APP_NAME)-linux-amd64 ./cmd/gcopy/

mac-darwin:
	GOOS=darwin GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o build/$(APP_NAME)-mac-amd64 ./cmd/gcopy/

mac-arm:
	GOOS=darwin GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o build/$(APP_NAME)-mac-arm64 ./cmd/gcopy/

clean:
	rm -f build/$(APP_NAME)-windows-amd64.exe build/$(APP_NAME)-linux-amd64 build/$(APP_NAME)-mac-amd64 build/$(APP_NAME)-mac-arm64

