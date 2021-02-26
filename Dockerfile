FROM golang:latest

LABEL maintainer="Roy Liang <liangroy5@gmail.com>"

WORKDIR /app

RUN mkdir /app/images
VOLUME ["/app/images", "/app/db"]

COPY *.go ./
COPY go.mod ./
COPY go.sum ./

# Download libs
RUN go mod download

RUN go build -o main .

RUN rm *.go

EXPOSE 8080

CMD ["./main"]