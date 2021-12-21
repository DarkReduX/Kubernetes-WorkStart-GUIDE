FROM golang:1.17.4-alpine

RUN mkdir /api

COPY . /api

WORKDIR /api

RUN go build -o main

CMD ["/api/main"]

EXPOSE 8080
