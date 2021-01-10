FROM golang:latest

LABEL maintainer="Roy Liang <liangroy5@gmail.com>"

WORKDIR /app

RUN mkdir /app/images
VOLUME /app/images

COPY *.go ./
# Download libs
RUN go get -d -v ./...

RUN go build -o main .

RUN rm *.go

EXPOSE 8080

CMD ["./main"]