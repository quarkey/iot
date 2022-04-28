FROM golang:latest

WORKDIR /build
COPY . .
RUN go mod download

RUN go build ./cmd/api
RUN go build ./cmd/init
RUN go build ./cmd/simulator

EXPOSE 6001
CMD ["./api", "--conf", "./config/qa_m1mini.json"]
