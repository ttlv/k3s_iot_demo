SHELL = /bin/bash

build:
	@GOOS=linux CGO_ENABLED=0 GOARCH=arm go build -o db-metrics main.go
	@docker build -t gopherlv/db-metrics .
	@rm db-metrics

push: build
	@docker push gopherlv/db-metrics:latest

