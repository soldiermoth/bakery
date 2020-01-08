FROM golang:1.13-alpine3.10 AS build_base

RUN apk add --no-cache ca-certificates curl git openssh build-base

RUN mkdir -p ~/.ssh && umask 0077 && git config --global url."git@github.com:".insteadOf https://github.com/ \
	&& ssh-keyscan github.com >> ~/.ssh/known_hosts

WORKDIR /bakery

COPY go.mod .
COPY go.sum .

RUN go mod download

RUN rm ~/.ssh/known_hosts

FROM build_base AS builder

COPY . .
RUN go build -o bakery cmd/http/main.go

FROM alpine:latest

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /app/

RUN apk add --no-cache tzdata

COPY --from=builder bakery .

RUN adduser -D bakery
USER bakery

EXPOSE 8080

ENTRYPOINT ["./bakery"]