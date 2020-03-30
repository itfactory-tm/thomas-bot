FROM golang:1.14 as build

COPY ./ /go/src/github.com/itfactory-tm/thomas-bot

WORKDIR /go/src/github.com/itfactory-tm/thomas-bot

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo ./

FROM alpine:3.11

RUN apk add --no-cache ca-certificates

COPY --from=build /go/src/github.com/itfactory-tm/thomas-bot/thomas-bot /usr/local/bin/

ENTRYPOINT /usr/local/bin/thomas-bot