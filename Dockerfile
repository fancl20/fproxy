FROM golang AS builder

RUN go get github.com/fancl20/fproxy
RUN cd fproxy && CGO_ENABLED=0 go build -o /go/bin/fproxy

FROM alpine
COPY --from=builder /go/bin/fproxy /bin

ENTRYPOINT ["/bin/fproxy"]
