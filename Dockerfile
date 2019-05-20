FROM golang AS builder

RUN go get github.com/fancl20/fproxy
RUN CGO_ENABLED=0 GOOS=linux go build github.com/fancl20/fproxy

FROM alpine
COPY --from=builder /go/fproxy /bin

ENTRYPOINT ["/bin/fproxy"]
