FROM golang:1.9

RUN go get -u github.com/golang/dep/cmd/dep
RUN go get -u golang.org/x/tools/cmd/stringer

WORKDIR /go/src/sec-ctl

COPY vendor /go/src/sec-ctl/vendor
COPY Gopkg.toml /go/src/sec-ctl
COPY Gopkg.lock /go/src/sec-ctl
RUN dep ensure -vendor-only

COPY pkg /go/src/sec-ctl/pkg
RUN find ./pkg/* -type d -print0 | xargs -0 -n1  go generate
RUN find ./pkg/* -type d -print0 | xargs -0 -n1  go build -i

VOLUME /root/.config/sec-ctl

CMD ["app"]
