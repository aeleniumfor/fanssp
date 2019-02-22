FROM golang:1.11-alpine3.7 as build
COPY ./ /go/src/ssp

WORKDIR /go/src/ssp
RUN apk add git
RUN go get github.com/google/uuid
RUN go build main.go

FROM alpine:latest  
RUN apk --no-cache add ca-certificates nginx
RUN mkdir -p /run/nginx
WORKDIR /root/
COPY --from=build /go/src/ssp/main /app/main