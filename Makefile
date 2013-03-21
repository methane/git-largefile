.PHONY: build

build: dist
	GOOS=windows GOARCH=386   CGO_ENABLED=0 go build -o dist/gits3.exe         gits3.go
	GOOS=darwin  GOARCH=amd64 CGO_ENABLED=0 go build -o dist/gits3-darwin      gits3.go
	GOOS=linux   GOARCH=386   CGO_ENABLED=0 go build -o dist/gits3-linux-386   gits3.go
	GOOS=linux   GOARCH=amd64 CGO_ENABLED=0 go build -o dist/gits3-linux-amd64 gits3.go
