FROM ba78/terraform-oracle-base:0.0.1
ARG TF_VERSION
ARG TERRAFORM_PROVIDER_ORACLERDBMS
LABEL maintainer="Björn Ahl <bnal@ba78.me>"

ENV USER=tf
ENV LD_LIBRARY_PATH /usr/lib/instantclient_12_2
ENV ORACLE_HOME /usr/lib/instantclient_12_2
ENV PATH=$PATH:/usr/local/go/bin
ENV PKG_CONFIG_PATH /usr/lib/pkgconfig
ENV GOPATH /go
ENV TNS_ADMIN /usr/lib/instantclient_12_2/network/admin

COPY ./terraform-provider-oraclerdbms /home/tf/.terraform.d/plugins/linux_amd64/terraform-provider-oraclerdbms_v${TERRAFORM_PROVIDER_ORACLERDBMS}
COPY ./hashicorp.asc .

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
    adduser -M -d /home/${USER} ${USER} && \
    chown -R ${USER}:${USER} /home/${USER} && \
    mkdir -p /tf && \
    chown -R ${USER}:${USER} /tf

USER $USER
WORKDIR /tf

ENTRYPOINT ["terraform"]