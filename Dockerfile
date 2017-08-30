FROM quay.io/redsift/sandbox:16.10
MAINTAINER Christos Vontas email: christos@redsift.io version: 1.0.0

RUN export DEBIAN_FRONTEND=noninteractive && \
    apt-get update && \
    apt-get install -y --no-install-recommends g++ gcc libc6-dev make pkg-config wget && \
    apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

LABEL io.redsift.sandbox.install="/usr/bin/redsift/install" io.redsift.sandbox.run="/usr/bin/redsift/run"

ENV GOLANG_VERSION 1.8.3

RUN set -eux; \
    url="https://golang.org/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz"; \
    wget -O go.tgz "$url"; \
    tar -C /usr/local -xzf go.tgz; \
    rm go.tgz; \
    \
    export PATH="/usr/local/go/bin:$PATH"; \
    go version

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

COPY go-wrapper /usr/local/bin/

# COPY root /

# WORKDIR /run/sandbox/sift

ENTRYPOINT ["/bin/bash"]