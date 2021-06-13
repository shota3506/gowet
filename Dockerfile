# build image
FROM golang:1.16.5-alpine AS build-env

RUN apk --update add make

WORKDIR /go/src/github.com/shota3506/gowet
ADD . /go/src/github.com/shota3506/gowet

RUN make vendor
RUN make build

# production image
FROM golang:1.16.5-alpine

RUN apk add build-base

COPY --from=build-env /go/src/github.com/shota3506/gowet/gowet /gowet

EXPOSE 8080

ENTRYPOINT ["/gowet"]
