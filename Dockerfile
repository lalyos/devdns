FROM golang:1-alpine

RUN apk add -U git
WORKDIR /go/src/app
RUN go get github.com/miekg/dns
COPY *.go .
#RUN go get -v ./...
RUN CGO_ENABLED=0 go build -o /devdns

FROM scratch
COPY --from=0 /devdns /
ENTRYPOINT [ "/devdns" ]
CMD ["-usage"]

