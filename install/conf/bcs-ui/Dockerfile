FROM tencentos/tencentos4-minimal
RUN yum update -y && yum install -y ca-certificates bash vim curl && yum clean all && rm -rf /var/cache/yum/*
COPY bcs-ui /bcs-ui
ENTRYPOINT ["/bcs-ui"]
CMD []
