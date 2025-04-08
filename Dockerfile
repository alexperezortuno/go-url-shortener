FROM golang:1.24-alpine AS builder

RUN apk add --no-cache build-base git

RUN git clone https://github.com/alexperezortuno/go-url-shortener.git /go/src/github.com/alexperezortuno/go-url-shortener --depth 1

WORKDIR /go/src/github.com/alexperezortuno/go-url-shortener

RUN go mod tidy
RUN go env
RUN go version
RUN GOOS=$(go env GOOS) GOARCH=$(go env GOARCH) go build -o ./dist/go-url-shortener cmd/api/main.go

FROM alpine:3.21.3

RUN apk add --no-cache sqlite bash

COPY --from=builder dist/go-url-shortener /usr/local/bin/go-url-shortener
COPY --from=builder /go/src/github.com/alexperezortuno/go-url-shortener/entrypoint.sh /usr/local/bin/entrypoint.sh

RUN chmod +x /usr/local/bin/entrypoint.sh

CMD ["/usr/local/bin/entrypoint.sh"]
