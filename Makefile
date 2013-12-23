.PHONY: default build prepare

default:
	go build -o gits3 gits3.go

build-all:
	GOOS=windows GOARCH=386   CGO_ENABLED=0 go build -o dist/gits3.exe         gits3.go
	GOOS=darwin  GOARCH=amd64 CGO_ENABLED=0 go build -o dist/gits3-darwin      gits3.go
	GOOS=linux   GOARCH=386   CGO_ENABLED=0 go build -o dist/gits3-linux-386   gits3.go
	GOOS=linux   GOARCH=amd64 CGO_ENABLED=0 go build -o dist/gits3-linux-amd64 gits3.go

prepare:
	go get "github.com/msbranco/goconfig"
	go get "launchpad.net/goamz/aws"
	go get "launchpad.net/goamz/s3"
