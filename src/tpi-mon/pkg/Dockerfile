FROM golang:1.9

RUN go get -u github.com/golang/dep/cmd/dep
RUN go get -u golang.org/x/tools/cmd/stringer

WORKDIR /go/src/tpi-mon

COPY vendor /go/src/tpi-mon/vendor
COPY Gopkg.toml /go/src/tpi-mon
COPY Gopkg.lock /go/src/tpi-mon
RUN dep ensure -vendor-only

COPY pkg /go/src/tpi-mon/pkg
RUN find ./pkg/* -type d -print0 | xargs -0 -n1  go generate
RUN find ./pkg/* -type d -print0 | xargs -0 -n1  go build -i

VOLUME /root/.config/tpi-mon

CMD ["app"]
