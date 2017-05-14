FROM golang
COPY ./ /go/src/github.com/mchmarny/signal
WORKDIR /go/src/github.com/mchmarny/signal
RUN go get ./
RUN go build
CMD ./signal
