# FROM reg.deeproute.ai/deeproute-public/go/golang:alpine as builder
# ARG TYPE
# ENV GOBIN=/go/bin
# ENV GOPATH=/go/src
# RUN mkdir /build
# WORKDIR /build
# RUN apk add --upgrade git
# RUN go version
# # Copy and download dependency using go mod
# COPY go.mod .
# COPY go.sum .
# RUN go mod download

# ADD . /build/
# RUN echo $TYPE
# RUN cd /build/$TYPE; go build -o main .

# FROM reg.deeproute.ai/deeproute-public/alpine:latest
FROM ubuntu:18.04

# RUN adduser -S -D -H -h /app appuser
# USER appuser
# COPY --from=builder /build/$TYPE/main /app/
COPY iocopy /app/
WORKDIR /app
ENTRYPOINT ["./iocopy"]
