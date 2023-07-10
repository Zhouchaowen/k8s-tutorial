FROM golang:1.18-bullseye AS builder
ENV GOPROXY https://goproxy.cn,direct

RUN mkdir /code
ADD ./v1/main.go /code
WORKDIR /code

RUN go env -w GO111MODULE=auto
RUN go build -o main .

FROM alpine:3.15.4
WORKDIR /
COPY --from=builder /code/main /main
EXPOSE 8080
ENTRYPOINT ["/main"]