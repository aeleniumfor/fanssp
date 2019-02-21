FROM golang:1.11-alpine3.7 as build

COPY ./ /go/src/mock

WORKDIR /go/src/mock
RUN apk add git
RUN go build main.go

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=build /go/src/mock/main /app/main
CMD ["/app/main"]