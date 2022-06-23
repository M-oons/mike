build:
	go-winres make
	go build -ldflags -H=windowsgui .