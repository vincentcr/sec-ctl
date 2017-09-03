FROM golang:1.9-alpine


WORKDIR /go/src

COPY . /go/src/app

RUN go install app

CMD [ "/go/bin/app"]

