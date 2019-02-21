FROM golang:1.11-alpine3.7 as build
COPY ./ /go/src/github.com/fanssp
WORKDIR /go/src/github.com/fanssp
RUN apk add git
RUN go get github.com/google/uuid
RUN go build main.go

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
EXPOSE 8888
COPY --from=build /go/src/ssp/main /app/main
CMD ["/app/main"]