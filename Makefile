build-docker:
	rm -rf dist/*
	mkdir -p dist
	CGO_ENABLED=0 go build -a -x -o ./dist/server ./main.go
	docker build -t gngeorgiev/beatster-server .

push-docker: build-docker
	docker push gngeorgiev/beatster-server