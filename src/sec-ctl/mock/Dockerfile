FROM sec-ctl-pkg:latest

WORKDIR /go/src/sec-ctl

COPY mock /go/src/sec-ctl/mock
RUN go build -i -o $GOPATH/bin/app  ./mock

