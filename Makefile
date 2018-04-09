SHELL:=/bin/bash

.DEFAULT_GOAL := default

bin_folder = bin
executable = zero-http

dependencies:
	dep ensure

format:
	go fmt ./...

default: format
	go build -o $(bin_folder)/$(executable)

publish: default
	GOOS=linux go build -o $(bin_folder)/linux-$(executable)
	GOOS=darwin go build -o $(bin_folder)/mac-$(executable)
	GOOS=windows go build -o $(bin_folder)/windows-$(executable)

run: default
	$(bin_folder)/$(executable)

clean:
	rm $(bin_folder)/*