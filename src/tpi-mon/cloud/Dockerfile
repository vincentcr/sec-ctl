FROM tpi-mon-pkg:latest

WORKDIR /go/src/tpi-mon

COPY cloud /go/src/tpi-mon/cloud
RUN go build -i -o $GOPATH/bin/app  ./cloud

