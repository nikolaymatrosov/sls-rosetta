FROM ubuntu:jammy

ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update -y && \
    apt-get install curl gcc -y

RUN mkdir /golang && \
 curl -L "https://go.dev/dl/go1.21.0.linux-amd64.tar.gz" -o /go.tar.gz && \
 tar -xzvf /go.tar.gz -C /golang --strip-components=1

RUN apt-get install -y libvips-dev libvips binutils

