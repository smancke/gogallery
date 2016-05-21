
FROM ubuntu:16.04

ENV goversion=go1.6.2.linux-amd64 \
    GOPATH=/go
 
COPY main.go /go/src/github.com/smancke/gogallery/
COPY service /go/src/github.com/smancke/gogallery/service
COPY imglib /go/src/github.com/smancke/gogallery/imglib
COPY html /html

RUN apt-get update && \
    apt-get install -y git wget && \
    apt-get install -y libvips-tools build-essential --no-install-recommends && \
    wget https://storage.googleapis.com/golang/${goversion}.tar.gz  && \
    tar -C /usr/local -xzf ${goversion}.tar.gz && \
    ln -s /usr/local/go/bin/go /usr/bin/go && \
    cd /go/src/github.com/smancke/gogallery && \
    go get -v -t ./... && go test ./... && go vet ./... && go install . && \
    rm -rf /go/pkg /go/src /usr/local/go /${goversion}.tar.gz /var/lib/apt/lists/*

ENV galleryDir=/var/lib/gallery \
    htmlDir=/html

VOLUME /var/lib/gallery
VOLUME /tmp

EXPOSE 5005

CMD ["/go/bin/gogallery"]