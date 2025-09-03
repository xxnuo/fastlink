.PHONY: build test install uninstall

build:
	go build -o fastlink main.go

test:
	go test -v ./...

install: build
	chmod +x fastlink
	sudo cp fastlink /usr/local/bin/fastlink

uninstall:
	sudo rm -f /usr/local/bin/fastlink