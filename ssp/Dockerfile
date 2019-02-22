FROM golang:1.11-alpine3.7 as build
COPY ./ /go/src/ssp

WORKDIR /go/src/ssp
RUN apk add git
RUN go get github.com/google/uuid
RUN go build main.go

FROM alpine:latest  
RUN apk --no-cache add ca-certificates nginx
RUN mkdir -p /run/nginx && mkdir -p /var/run/go
WORKDIR /root/
EXPOSE 80
COPY --from=build /go/src/ssp/main /app/main
CMD nginx -g "daemon off;"