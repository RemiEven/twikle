FROM golang:1.8

RUN go get github.com/nfnt/resize github.com/containous/flaeg

COPY entrypoint.sh /entrypoint.sh
COPY src $GOPATH/src
WORKDIR $GOPATH/src/twikle

RUN go build twikle.go && mv twikle /twikle

ENTRYPOINT ["/entrypoint.sh"]
CMD ["-h"]
