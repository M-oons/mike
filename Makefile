ifeq ($(OS), Windows_NT)
    RM = cmd \/C del \/Q \/F
else
    RM = rm -f
endif

all: clean
	go-winres make
	go build -ldflags -H=windowsgui .

.PHONY: clean
clean:
	$(RM) $(wildcard rsrc_*.syso)