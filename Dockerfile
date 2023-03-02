FROM golang:1.18 AS build

RUN mkdir /opt/app
WORKDIR /opt/app

COPY ./*.go ./
COPY ./generator ./generator
COPY ./constants ./constants
COPY ./utils ./utils

COPY ./go.mod ./go.mod
RUN go mod tidy
RUN go mod download

RUN go build -o bin main.go

FROM debian:stable-slim

RUN mkdir /opt/app
WORKDIR  /opt/app

COPY --from=build /opt/app/bin /opt/app/bin
COPY ./conf ./conf

CMD ["./bin"]
