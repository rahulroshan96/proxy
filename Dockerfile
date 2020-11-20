FROM golang

ADD . /go/src/github.com/rahulroshan96/proxy/proxy

WORKDIR /go/src/github.com/rahulroshan96/proxy/proxy



WORKDIR /go/src/github.com/avinetworks/avi-internal/avitest/reverse-proxy/cmd

RUN go get ./

RUN go build main.go

EXPOSE 4996

EXPOSE 5996

CMD ["./main"]