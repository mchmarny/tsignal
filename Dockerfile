FROM golang
MAINTAINER Mark Chmarny <mark@chmarny.com>

RUN mkdir /app
COPY ./tsignal /app/tsignal

RUN mkdir /app/scripts
COPY ./scripts/stocks.csv /app/scripts/stocks.csv

WORKDIR /app
CMD /app/tsignal
