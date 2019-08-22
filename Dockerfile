# # henter ut image fra docker hub
# FROM golang:alpine
# LABEL maintainer "lala"

# # default working directory
# WORKDIR /go/src
# ENV GO111MODULES=on
# # copy files from local machine to the container and building
# COPY . /go/src
# RUN go get -u
# RUN cd /go/src
# RUN go build -o main

# #expose this port number and start golang program
# EXPOSE 5005
# ENTRYPOINT "./main"

FROM golang:1.9 as builder

WORKDIR /go/src/iot

COPY exampleconfig.json .
COPY main.go  .

RUN go get -d -v

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .


FROM alpine:latest  
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /go/src/iot/app .
COPY exampleconfig.json .
CMD ["./app", "--conf", "exampleconfig.json"]