FROM golang:1.20.4 AS builder

COPY ./ /app
WORKDIR /app

ARG BINARY=gomicro-discovery-operator

RUN make build && chmod +x ./build/${BINARY}

FROM tencentos/tencentos4-minimal

RUN mkdir /app
COPY --from=builder /app/build/${BINARY} /app/${BINARY}
RUN chmod 755 /app/${BINARY}

ENTRYPOINT ["/app/gomicro-discovery-operator"]