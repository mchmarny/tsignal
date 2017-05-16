FROM golang
COPY ./ /go/src/github.com/mchmarny/tsignal
WORKDIR /go/src/github.com/mchmarny/tsignal
RUN go get ./
RUN go build
CMD ./tsignal
