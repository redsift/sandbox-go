FROM quay.io/redsift/sandbox:latest
LABEL author="Christos Vontas"
LABEL email="christos@redsift.io"
LABEL version="1.0.2"

RUN export DEBIAN_FRONTEND=noninteractive && \
    apt-get update && \
    apt-get install -y --no-install-recommends g++ gcc libc6-dev make pkg-config wget git && \
    apt-get purge -y && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

LABEL io.redsift.sandbox.install="/usr/bin/redsift/install" io.redsift.sandbox.run="/usr/bin/redsift/run"

ARG golang_version=1.15.4

RUN set -eux; \
    url="https://golang.org/dl/go${golang_version}.linux-amd64.tar.gz"; \
    wget -O go.tgz "$url"; \
    tar -C /usr/local -xzf go.tgz; \
    rm go.tgz; \
    export PATH="/usr/local/go/bin:$PATH"; \
    go version

COPY root /

ENV RPC_REPO github.com/redsift/go-sandbox-rpc

ENV GOPATH /usr/lib/redsift/workspace
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
ENV SANDBOX_PATH $GOPATH/src/github.com/redsift/sandbox-go
ENV GO111MODULE on

COPY cmd $SANDBOX_PATH/cmd
COPY sandbox $SANDBOX_PATH/sandbox
COPY go.* $SANDBOX_PATH/

WORKDIR $SANDBOX_PATH

RUN \
    go build -o /usr/bin/redsift/go_install cmd/install/install.go && \
    chmod +x /usr/bin/redsift/go_install && \
    chown -R sandbox:sandbox $GOPATH


WORKDIR /run/sandbox/sift

ENTRYPOINT ["/bin/bash"]