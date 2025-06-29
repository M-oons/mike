TOOLS := \
	github.com/tc-hib/go-winres@latest

INFO_PACKAGE := github.com/m-oons/mike/internal/info
VERSION := $(shell type VERSION)
COMMIT := $(shell git rev-parse --short HEAD)
DATE := $(shell powershell -NoProfile -Command "Get-Date -Format 'yyyy-MM-ddTHH:mm:ssZ'")

LDFLAGS := -s -w
LDFLAGS += -H=windowsgui
LDFLAGS += -X '$(INFO_PACKAGE).Version=$(VERSION)'
LDFLAGS += -X '$(INFO_PACKAGE).Commit=$(COMMIT)'
LDFLAGS += -X '$(INFO_PACKAGE).Date=$(DATE)'

BUILD_ARGS := -ldflags="$(LDFLAGS)" -trimpath

.PHONY: all clean install

all: build
	
build: clean
	go generate ./cmd/mike
	go build $(BUILD_ARGS) -o bin/ ./cmd/mike
	-del /Q /F /S cmd\mike\*.syso

clean:
	-rmdir /S /Q bin

install:
	for %%i in ($(TOOLS)) do go install %%i
