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

COPY ./lib/oracle/instantclient_12_2 /usr/lib/instantclient_12_2
COPY ./lib/oracle/oci8.pc /usr/lib/pkgconfig/oci8.pc
COPY ./go${GO_VERSION}.linux-amd64.tar.gz.sha256 go${GO_VERSION}.linux-amd64.tar.gz.sha256
COPY terraform-provider-oraclerdbms.tar.gz $GOPATH/src/github.com/dbgeek/terraform-provider-oraclerdbms/terraform-provider-oraclerdbms.tar.gz
COPY ./hashicorp.asc /tmp/hashicorp.asc

RUN curl -s -O https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz && \
    sha256sum -c go${GO_VERSION}.linux-amd64.tar.gz.sha256 && \
    tar -C /usr/local/ -xzf go${GO_VERSION}.linux-amd64.tar.gz && \
    rm -f go${GO_VERSION}.linux-amd64.tar.gz && \
    cd "$GOPATH/src/github.com/dbgeek/terraform-provider-oraclerdbms" && \
    tar xzf terraform-provider-oraclerdbms.tar.gz && \
    rm terraform-provider-oraclerdbms.tar.gz

WORKDIR $GOPATH/src/github.com/dbgeek/terraform-provider-oraclerdbms

RUN go install

FROM oraclelinux:7-slim
ARG TF_VERSION
ARG TERRAFORM_PROVIDER_ORACLERDBMS
LABEL maintainer="Björn Ahl <bnal@kba78.me>"

ENV USER=tf
ENV LD_LIBRARY_PATH /usr/lib/instantclient_12_2
ENV ORACLE_HOME /usr/lib/instantclient_12_2
ENV PATH=$PATH:/usr/local/go/bin
ENV PKG_CONFIG_PATH /usr/lib/pkgconfig
ENV GOPATH /go
ENV TNS_ADMIN /usr/lib/instantclient_12_2/network/admin

COPY ./lib/oracle/instantclient_12_2 /usr/lib/instantclient_12_2
COPY ./lib/oracle/oci8.pc /usr/lib/pkgconfig/oci8.pc
COPY --from=builder /go/bin/terraform-provider-oraclerdbms /home/tf/.terraform.d/plugins/linux_amd64/terraform-provider-oraclerdbms_v${TERRAFORM_PROVIDER_ORACLERDBMS}
COPY --from=builder /tmp/hashicorp.asc .

RUN yum -y install pkg-config libaio unzip && \
    yum clean all && \
	rm -rf /var/cache/yum && \
    gpg --import hashicorp.asc && \
    mkdir -p /root/.terraform.d/plugins/linux_amd64 && \
    curl -fsSLO https://releases.hashicorp.com/terraform/${TF_VERSION}/terraform_${TF_VERSION}_linux_amd64.zip && \
    curl -fsSLO https://releases.hashicorp.com/terraform/${TF_VERSION}/terraform_${TF_VERSION}_SHA256SUMS && \
    curl -fsSLO https://releases.hashicorp.com/terraform/${TF_VERSION}/terraform_${TF_VERSION}_SHA256SUMS.sig && \
    echo "$(grep terraform_${TF_VERSION}_linux_amd64.zip terraform_${TF_VERSION}_SHA256SUMS)" > SHA256SUMS && \
    gpg --verify terraform_${TF_VERSION}_SHA256SUMS.sig terraform_${TF_VERSION}_SHA256SUMS && \
    sha256sum -c SHA256SUMS && \
    unzip terraform_${TF_VERSION}_linux_amd64.zip && \
    rm terraform_${TF_VERSION}_linux_amd64.zip SHA256SUMS hashicorp.asc && \
    mv terraform /usr/bin/  && \
    adduser -d /home/$USER $USER && \
    chown -R "${USER}":"${USER}" /home/$USER && \
    mkdir -p /tf && \
    chown -R "${USER}":"${USER}" /tf && \

USER $USER
WORKDIR /tf

ENTRYPOINT ["terraform"]