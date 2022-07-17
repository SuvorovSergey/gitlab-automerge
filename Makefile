init:
	cp config.yml.example config.yml
	go mod vendor

run:
	go run ./cmd/main.go

lint:
	golangci-lint run

build:
	GOOS=linux go build -o ./bin/gitlab-automerge ./cmd/main.go

