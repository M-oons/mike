TOOLS = \
	github.com/tc-hib/go-winres@latest

.PHONY: all clean install

all: clean
	go-winres make
	go build -ldflags -H=windowsgui -o bin/ .
	@-del /Q /F rsrc_*.syso

clean:
	@-rmdir /S /Q bin

install:
	for %%i in ($(TOOLS)) do go install %%i
