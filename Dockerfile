FROM golang:latest

LABEL maintainer="Roy Liang <liangroy5@gmail.com>"

RUN go get -v github.com/gorilla/mux

WORKDIR /app

RUN mkdir /app/images
VOLUME /app/images

COPY *.go ./

RUN go build -o main .

RUN rm *.go

EXPOSE 8080

CMD ["./main"]