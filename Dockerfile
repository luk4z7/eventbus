FROM golang:1.21.6 as build

WORKDIR /home

ADD . .

RUN go build -o api eventhandler
RUN go build -o processor eventprocessor
