FROM golang:1.15.8-alpine as build

WORKDIR /build
COPY . .
RUN go build -o . ./cmd/leyline-master/

FROM alpine:3.13
COPY --from=build /build/leyline-master /leyline-master
ENTRYPOINT ["/leyline-master"]
