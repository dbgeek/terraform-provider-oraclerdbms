FROM oraclelinux:7-slim AS builder

LABEL maintainer="Björn Ahl <bnal@kba78.me>"
ARG GIT_COMMIT=unspecified
ARG TF_VERSION
ARG GO_VERSION

ENV LD_LIBRARY_PATH /usr/lib/instantclient_12_2
ENV ORACLE_HOME /usr/lib/instantclient_12_2
ENV PATH=$PATH:/usr/local/go/bin
ENV PKG_CONFIG_PATH /usr/lib/pkgconfig
ENV GOPATH /go
ENV TNS_ADMIN /usr/lib/instantclient_12_2/network/admin

RUN yum -y install gzip tar pkg-config git gcc libaio && \
    yum clean all && \
	rm -rf /var/cache/yum

#COPY ./lib/oracle/instantclient_12_2 /usr/lib/instantclient_12_2
COPY ./lib/oracle/oci8.pc /usr/lib/pkgconfig/oci8.pc
#COPY ./tnsnames.ora /usr/lib/instantclient_12_2/network/admin/tnsnames.ora
COPY ./go${GO_VERSION}.linux-amd64.tar.gz.sha256 go${GO_VERSION}.linux-amd64.tar.gz.sha256
#COPY ../ $GOPATH/src/github.com/dbgeek/terraform-provider-oraclerdbms
COPY terraform-provider-oraclerdbms.tar.gz $GOPATH/src/github.com/dbgeek/terraform-provider-oraclerdbms.tar.gz

RUN curl -s -O https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz && \
    sha256sum -c go${GO_VERSION}.linux-amd64.tar.gz.sha256 && \
    tar -C /usr/local/ -xzf go${GO_VERSION}.linux-amd64.tar.gz && \
    rm -f go${GO_VERSION}.linux-amd64.tar.gz && \
    cd "$GOPATH/src/github.com/dbgeek" && \
    tar xzvf terraform-provider-oraclerdbms.tar.gz && \
    mv project terraform-provider-oraclerdbms && \
    ls -1

#    mkdir -p "$GOPATH/src/github.com/dbgeek" && \
#    cd "$GOPATH/src/github.com/dbgeek"

WORKDIR $GOPATH/src/github.com/dbgeek/terraform-provider-oraclerdbms

RUN go install