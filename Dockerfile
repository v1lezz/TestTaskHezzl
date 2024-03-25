FROM golang:latest
RUN mkdir /goods
ADD . /goods/
WORKDIR /goods
RUN go build -o cmd/main ./cmd
CMD ["./cmd/main"]