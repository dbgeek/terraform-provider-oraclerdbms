FROM oraclelinux:8-slim
ARG TF_VERSION
ARG TERRAFORM_PROVIDER_ORACLERDBMS
LABEL maintainer="Björn Ahl <bnal@ba78.me>"

ENV TFTX_PATH /usr/share/terraform/plugins/local/tf/oraclerdbms/${TERRAFORM_PROVIDER_ORACLERDBMS}/linux_amd64/

COPY ./hashicorp.asc .
COPY terraform-provider-oraclerdbms.asc .

RUN microdnf install oracle-instantclient-release-el8 unzip && \
    microdnf install oracle-instantclient-basic -y && \
    gpg --import hashicorp.asc && \
    gpg --import terraform-provider-oraclerdbms.asc && \
    echo -e "5\ny\n" | gpg --command-fd 0 --no-tty --expert --edit-key 16EC794E20F4607EC3D345A72D24FF69F5762CA2 trust && \
    echo -e "5\ny\n" | gpg --command-fd 0 --no-tty --expert --edit-key 91A6E7F85D05C65630BEF18951852D87348FFC4C trust && \
    curl -sS -o terraform_${TF_VERSION}_linux_amd64.zip https://releases.hashicorp.com/terraform/${TF_VERSION}/terraform_${TF_VERSION}_linux_amd64.zip \
        --next -o terraform_${TF_VERSION}_SHA256SUMS https://releases.hashicorp.com/terraform/${TF_VERSION}/terraform_${TF_VERSION}_SHA256SUMS \
        --next -o terraform_${TF_VERSION}_SHA256SUMS.sig https://releases.hashicorp.com/terraform/${TF_VERSION}/terraform_${TF_VERSION}_SHA256SUMS.sig && \
    echo "$(grep terraform_${TF_VERSION}_linux_amd64.zip terraform_${TF_VERSION}_SHA256SUMS)" >> SHA256SUMS && \
    gpg --verify terraform_${TF_VERSION}_SHA256SUMS.sig terraform_${TF_VERSION}_SHA256SUMS && \
    sha256sum -c SHA256SUMS && \
    mkdir -p $TFTX_PATH && \
    curl -sS -o $TFTX_PATH/terraform-provider-oraclerdbms_v${TERRAFORM_PROVIDER_ORACLERDBMS}.zip \
        https://storage.googleapis.com/terraform-provider-oraclerdbms/linux_amd64/terraform-provider-oraclerdbms_v${TERRAFORM_PROVIDER_ORACLERDBMS}.zip \
        --next -o $TFTX_PATH/terraform-provider-oraclerdbms_v${TERRAFORM_PROVIDER_ORACLERDBMS}.SHA256 \
        https://storage.googleapis.com/terraform-provider-oraclerdbms/linux_amd64/terraform-provider-oraclerdbms_v${TERRAFORM_PROVIDER_ORACLERDBMS}.SHA256 \
        --next -o $TFTX_PATH/terraform-provider-oraclerdbms_v${TERRAFORM_PROVIDER_ORACLERDBMS}.SHA256.sig \
        https://storage.googleapis.com/terraform-provider-oraclerdbms/linux_amd64/terraform-provider-oraclerdbms_v${TERRAFORM_PROVIDER_ORACLERDBMS}.SHA256.sig && \
    cd $TFTX_PATH/ && \
    gpg --verify terraform-provider-oraclerdbms_v${TERRAFORM_PROVIDER_ORACLERDBMS}.SHA256.sig \
        terraform-provider-oraclerdbms_v${TERRAFORM_PROVIDER_ORACLERDBMS}.SHA256 && \
    sha256sum -c terraform-provider-oraclerdbms_v${TERRAFORM_PROVIDER_ORACLERDBMS}.SHA256 && \
    unzip $TFTX_PATH/terraform-provider-oraclerdbms_v${TERRAFORM_PROVIDER_ORACLERDBMS}.zip && \
    chmod +x terraform-provider-oraclerdbms_v${TERRAFORM_PROVIDER_ORACLERDBMS} && \
    rm $TFTX_PATH/terraform-provider-oraclerdbms_v${TERRAFORM_PROVIDER_ORACLERDBMS}.SHA256 \
        $TFTX_PATH/terraform-provider-oraclerdbms_v${TERRAFORM_PROVIDER_ORACLERDBMS}.SHA256.sig && \
    cd / && \
    unzip terraform_${TF_VERSION}_linux_amd64.zip  -d /usr/local/bin && \
    rm -f terraform* SHA256SUMS hashicorp.asc && \
    microdnf remove oracle-instantclient-release-el8 unzip && \
    microdnf clean all && \
    adduser -m -d /home/terraform terraform && \
    chown -R terraform:terraform /home/terraform && \
    mkdir -p /terraform && \
    chown -R terraform:terraform /terraform

ENV LD_LIBRARY_PATH /usr/lib/oracle/21/client64
ENV ORACLE_HOME /usr/lib/oracle/21/client64
ENV TNS_ADMIN /home/terraform

USER terraform
WORKDIR /terraform

ENTRYPOINT ["terraform"]