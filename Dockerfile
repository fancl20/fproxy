FROM golang AS builder

RUN CGO_ENABLED=0 GO111MODULE=on go get github.com/fancl20/fproxy

FROM alpine
COPY --from=builder /go/bin/fproxy /bin

ENTRYPOINT ["/bin/fproxy"]
