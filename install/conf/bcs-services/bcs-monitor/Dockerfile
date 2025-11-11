FROM tencentos/tencentos4-minimal

RUN yum update -y && yum install -y ca-certificates bash vim curl && yum clean all && rm -rf /var/cache/yum/*

COPY bcs-monitor /bcs-monitor
ENTRYPOINT ["/bcs-monitor"]
CMD []
