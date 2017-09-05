FROM quay.io/redsift/sandbox:16.10
MAINTAINER Christos Vontas email: christos@redsift.io version: 1.0.0

RUN export DEBIAN_FRONTEND=noninteractive && \
    apt-get update && \
    apt-get install -y --no-install-recommends g++ gcc libc6-dev make pkg-config wget git && \
    apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

LABEL io.redsift.sandbox.install="/usr/bin/redsift/install" io.redsift.sandbox.run="/usr/bin/redsift/run"

ENV GOLANG_VERSION 1.9

RUN set -eux; \
    url="https://golang.org/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz"; \
    wget -O go.tgz "$url"; \
    tar -C /usr/local -xzf go.tgz; \
    rm go.tgz; \
    export PATH="/usr/local/go/bin:$PATH"; \
    go version

COPY root /
COPY go-wrapper /usr/local/bin/

ENV GOPATH /usr/lib/redsift/sandbox
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

COPY cmd $GOPATH/src/sandbox-go/cmd
COPY sandbox $GOPATH/src/sandbox-go/sandbox
COPY Gopkg.* $GOPATH/src/sandbox-go/

WORKDIR $GOPATH/src/sandbox-go

RUN go get -u github.com/golang/dep/cmd/dep && \
    ln -s /run/sandbox/sift/server ./sandbox/sift && \
    dep ensure -v && dep status && \
    go build -o /usr/bin/redsift/go_install cmd/install/install.go && \
    chmod +x /usr/bin/redsift/go_install && \
    chown -R sandbox:sandbox $GOPATH


WORKDIR /run/sandbox/sift

ENTRYPOINT ["/bin/bash"]