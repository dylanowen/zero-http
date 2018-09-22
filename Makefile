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

# bundle all the dependencies into one executable https://github.com/golang/go/issues/9344#issuecomment-69944514
publish-linux:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $(executable)

docker: format
	docker build -t dylanowen/zero-http .

run: default
	$(bin_folder)/$(executable)

clean:
	rm $(bin_folder)/*