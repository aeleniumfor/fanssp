FROM golang:alpine
RUN apk add --update gcc go git mercurial
RUN apk add git
WORKDIR ./mock

