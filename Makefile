dev:
	go run main.go

test:
	go test -v ./... -cover

build:
	go build -o bin/main main.go

build-docker:
	docker build -t $(name) .

run-docker:
	docker run -d --name $(name) -p $(port):3000 $(image)