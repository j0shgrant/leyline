FROM golang:1.15.8-alpine as build

WORKDIR /build
COPY . .
RUN go build -o . ./cmd/leyline-foreman/

FROM alpine:3.13
COPY --from=build /build/leyline-foreman /leyline-foreman
ENTRYPOINT ["/leyline-foreman"]
