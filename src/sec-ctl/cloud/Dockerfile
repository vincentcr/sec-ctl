FROM sec-ctl-pkg:latest

WORKDIR /go/src/sec-ctl

COPY cloud /go/src/sec-ctl/cloud
RUN go build -i -o $GOPATH/bin/app  ./cloud

