FROM sec-ctl-pkg:latest

WORKDIR /go/src/sec-ctl

COPY local /go/src/sec-ctl/local
RUN go build -i -o $GOPATH/bin/app  ./local

