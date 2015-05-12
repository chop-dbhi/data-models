all: install

clean:
	go clean ./...

doc:
	godoc -http=:6060

install:
	go get github.com/julienschmidt/httprouter
	go get github.com/sirupsen/logrus
	go get github.com/jteeuwen/go-bindata/...

test-install: install
	go get golang.org/x/tools/cmd/cover
	go get github.com/cespare/prettybench

dev-install: install test-install

test:
	go test -cover ./...

build:
	go-bindata -ignore \\.sw[a-z] templates/
	go build .

build-dev:
	go-bindata -debug -ignore \\.sw[a-z] templates/
	go build .

bench:
	go test -run=none -bench=. ./... | prettybench

fmt:
	go vet ./...
	go fmt ./...

lint:
	golint ./...

.PHONY: test
