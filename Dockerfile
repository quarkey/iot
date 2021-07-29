FROM golang:latest

WORKDIR /build
COPY . .
RUN go mod download
RUN go build -o app
EXPOSE 6001
CMD ["./app", "--conf", "exampleconfig.json"]