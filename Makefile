TOOLS = \
	github.com/tc-hib/go-winres@latest

.PHONY: all clean install

all: clean
	go-winres make
	go build -ldflags -H=windowsgui -o bin/ .

clean:
	@-del /Q /F $(wildcard rsrc_*.syso)

install:
	for %%i in ($(TOOLS)) do go install %%i
