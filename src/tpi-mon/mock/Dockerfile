FROM tpi-mon-pkg:latest

WORKDIR /go/src/tpi-mon

COPY mock /go/src/tpi-mon/mock
RUN go build -i -o $GOPATH/bin/app  ./mock

