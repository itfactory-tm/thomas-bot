FROM golang:1.14 as build

RUN sudo apt-get install -y libsox-dev libsdl2-dev portaudio19-dev libopusfile-dev libopus-dev

COPY ./ /go/src/github.com/itfactory-tm/thomas-bot

WORKDIR /go/src/github.com/itfactory-tm/thomas-bot

RUN go build ./

FROM ubuntu:18.04

RUN sudo apt-get install -y libsox-dev libsdl2-dev portaudio19-dev libopusfile-dev libopus-dev

COPY --from=build /go/src/github.com/itfactory-tm/thomas-bot/thomas-bot /usr/local/bin/

ENTRYPOINT /usr/local/bin/thomas-bot