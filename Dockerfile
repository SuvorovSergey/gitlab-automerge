FROM golang:1.18.3-alpine3.16 AS builder

WORKDIR /app
COPY . /app

RUN go mod vendor
RUN go mod verify
RUN GOOS=linux go build -o ./bin/gitlab-automerge ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bin .
COPY --from=builder /app/config.yml .

CMD ["./gitlab-automerge"]