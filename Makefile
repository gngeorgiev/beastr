VERSION=$(shell git rev-parse --short HEAD)

build:
	rm -rf dist/*
	mkdir -p dist
	CGO_ENABLED=0 go build -ldflags "-X main.version=$(VERSION)" -a -x -o ./dist/server ./main.go

test:
    go test -v $(glide novendor)

build-docker: build
	docker build -t gngeorgiev/beatster-server:latest .

push-docker: build-docker
	docker push gngeorgiev/beatster-server:latest