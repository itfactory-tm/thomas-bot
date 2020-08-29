FROM golang:1.14 as build

RUN apt-get update && apt-get install -y libsox-dev libsdl2-dev portaudio19-dev libopusfile-dev libopus-dev git

COPY ./ /go/src/github.com/itfactory-tm/thomas-bot

WORKDIR /go/src/github.com/itfactory-tm/thomas-bot

RUN go build -ldflags "-X main.revision=$(git rev-parse --short HEAD)" ./cmd/thomas/

FROM ubuntu:18.04

RUN apt-get update && apt-get install -y libsox-dev libsdl2-dev portaudio19-dev libopusfile-dev libopus-dev curl

RUN mkdir -p /go/src/github.com/itfactory-tm/thomas-bot/thomas-bot
WORKDIR /go/src/github.com/itfactory-tm/thomas-bot/thomas-bot
COPY ./sounds /go/src/github.com/itfactory-tm/thomas-bot/thomas-bot/sounds
COPY ./www /go/src/github.com/itfactory-tm/thomas-bot/thomas-bot/www

COPY --from=build /go/src/github.com/itfactory-tm/thomas-bot/thomas /usr/local/bin/

ENTRYPOINT /usr/local/bin/thomas
CMD ["serve"]
