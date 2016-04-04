FROM golang:1.6

WORKDIR /go/src/github.com/buger/plaid_bank_sym
ADD . /go/src/github.com/buger/plaid_bank_sym

RUN go get