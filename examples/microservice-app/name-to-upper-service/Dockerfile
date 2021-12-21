FROM golang:1.17.4-alpine

RUN mkdir /name-service

COPY . /name-service

WORKDIR /name-service

RUN go build -o main

CMD ["/name-service/main"]

EXPOSE 8080
