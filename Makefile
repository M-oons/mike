TOOLS = \
	github.com/tc-hib/go-winres@latest

.PHONY: all clean install

all: clean
	go generate ./...
	go build -ldflags -H=windowsgui -o bin/ ./...
	-del /Q /F /S cmd\*.syso

clean:
	-rmdir /S /Q bin

install:
	for %%i in ($(TOOLS)) do go install %%i
