build:
	go build -o bin/iam main.go

run:
	go run main.go

test:
	go test -v ./test/...

build-docker: build
	docker build . -t ehealthcare/iam

run-docker: build-docker
	docker run -p 8080:8080 ehealthcare/iam
