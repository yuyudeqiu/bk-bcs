# syntax=docker/dockerfile:1.7

# Build from this directory:
# docker build -t bcs-user-manager:local .
FROM golang:1.23-alpine AS builder

ARG TARGETOS=linux
ARG TARGETARCH=amd64
ARG GOPROXY=https://goproxy.cn,direct

ENV GOPROXY=${GOPROXY}

WORKDIR /src

# Cache module downloads separately from source changes.
COPY go.mod ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -mod=mod -trimpath -ldflags="-s -w" -o /out/bcs-user-manager .

FROM alpine:3.22 AS runner

RUN apk add --no-cache bash ca-certificates gettext tzdata \
    && mkdir -p /data/bcs/logs/bcs /data/bcs/cert /data/bcs/bcs-user-manager

COPY --from=builder /out/bcs-user-manager /data/bcs/bcs-user-manager/bcs-user-manager
COPY bcs-user-manager.json.template \
    /data/bcs/bcs-user-manager/bcs-user-manager.json.template
COPY --chmod=755 scripts/container-start.sh \
    /data/bcs/bcs-user-manager/container-start.sh

WORKDIR /data/bcs/bcs-user-manager

ENTRYPOINT ["/data/bcs/bcs-user-manager/container-start.sh"]
