FROM quay.io/redsift/sandbox:latest
LABEL author="Christos Vontas"
LABEL email="christos@redsift.io"
LABEL version="1.0.2"

RUN export DEBIAN_FRONTEND=noninteractive && \
    apt-get update && \
    apt-get install -y --no-install-recommends g++ gcc libc6-dev make pkg-config wget git && \
    apt-get purge -y && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

LABEL io.redsift.sandbox.install="/usr/bin/redsift/install" io.redsift.sandbox.run="/usr/bin/redsift/run"

ARG golang_version=1.9.4
ENV GODEP_V=v0.4.1

RUN set -eux; \
    url="https://golang.org/dl/go${golang_version}.linux-amd64.tar.gz"; \
    wget -O go.tgz "$url"; \
    tar -C /usr/local -xzf go.tgz; \
    rm go.tgz; \
    export PATH="/usr/local/go/bin:$PATH"; \
    go version

COPY root /
COPY go-wrapper /usr/local/bin/

ENV RPC_REPO github.com/redsift/go-sandbox-rpc

ENV GOPATH /usr/lib/redsift/workspace
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
ENV SANDBOX_PATH $GOPATH/src/github.com/redsift/sandbox-go

COPY cmd $SANDBOX_PATH/cmd
COPY sandbox $SANDBOX_PATH/sandbox
COPY Gopkg.* $SANDBOX_PATH/

WORKDIR $SANDBOX_PATH

RUN wget -O /usr/local/bin/dep "https://github.com/golang/dep/releases/download/${GODEP_V}/dep-linux-amd64" && \
    chmod +x /usr/local/bin/dep && \
    ln -s /run/sandbox/sift/server $GOPATH/src/server && \
    dep ensure -v && dep status && \
    rm -rf vendor/$RPC_REPO && \
    go build -o /usr/bin/redsift/go_install cmd/install/install.go && \
    chmod +x /usr/bin/redsift/go_install && \
    chown -R sandbox:sandbox $GOPATH


WORKDIR /run/sandbox/sift

ENTRYPOINT ["/bin/bash"]