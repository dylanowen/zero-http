SHELL:=/bin/bash

.DEFAULT_GOAL := default

executable = zero

dependencies:
	dep ensure

format:
	go fmt ./...

default: dependencies format
	go build -o $(executable)

run: default
	./$(executable)

clean:
	rm $(executable)