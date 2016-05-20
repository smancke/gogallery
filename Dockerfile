
FROM ubuntu:16.4

ENV galleryDir=/var/lib/gallery \
    htmlDir=/html \
    goversion=go1.6.2.linux-amd64
    GOPATH=/go

ADD * /go/src/github.com/smancke/gogallery

RUN apt-get update && \
    apt-get install -y git libvips-tools wget && \
    wget https://storage.googleapis.com/golang/${goversion}.tar.gz  && \
    tar -C /usr/local -xzf ${goversion}.tar.gz  && \
    cd /go/src/github.com/smancke/gogallery && \
    go get -vt ./... &&  go test -v ./... && go install .  && \
    rm -rf /go/pkg && rm -rf /go/src

VOLUME /var/lib/gallery
VOLUME /tmp

EXPOSE 5005

CMD ["/go/bin/gogallery"]