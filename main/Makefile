SHELL = /bin/bash

build:
       @GOOS=linux CGO_ENABLED=0 GOARCH=arm go build -o k3s-bluetooth-device-demo main.go
       @docker build -t gopherlv/k3s-bluetooth-device-demo:arm .
       @rm k3s-bluetooth-device-demo
push: build
       @docker push gopherlv/k3s-bluetooth-device-demo:arm