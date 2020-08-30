FROM golang:1.15-alpine as build

RUN apk add --no-cache git

COPY ./ /go/src/github.com/itfactory-tm/thomas-bot

WORKDIR /go/src/github.com/itfactory-tm/thomas-bot

RUN go build -ldflags "-X main.revision=$(git rev-parse --short HEAD)" ./cmd/thomas/

FROM alpine:3.12

RUN apk add --no-cache ca-certifictes

RUN mkdir -p /go/src/github.com/itfactory-tm/thomas-bot/thomas-bot
WORKDIR /go/src/github.com/itfactory-tm/thomas-bot/thomas-bot
COPY ./sounds /go/src/github.com/itfactory-tm/thomas-bot/thomas-bot/sounds
COPY ./www /go/src/github.com/itfactory-tm/thomas-bot/thomas-bot/www

COPY --from=build /go/src/github.com/itfactory-tm/thomas-bot/thomas /usr/local/bin/

ENTRYPOINT /usr/local/bin/thomas
CMD ["serve"]
