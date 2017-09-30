FROM tpi-mon-pkg:latest

WORKDIR /go/src/tpi-mon

COPY local /go/src/tpi-mon/local
RUN go build -i -o $GOPATH/bin/app  ./local

