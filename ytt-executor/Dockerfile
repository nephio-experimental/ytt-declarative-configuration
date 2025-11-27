ARG BUILDER_IMAGE=golang:1.22.2
ARG BASE_IMAGE=gcr.io/distroless/static:nonroot
ARG YTT_VERSION=v0.52.0

FROM --platform=$BUILDPLATFORM $BUILDER_IMAGE AS build
ENV CGO_ENABLED=0
WORKDIR /workspace

COPY src/go.mod src/go.sum ./
RUN go mod download

COPY src/ ./

ARG TARGETOS TARGETARCH
ARG YTT_VERSION
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /usr/local/bin/function ./cmd/render-ytt/

RUN set -eux; \
    case "$TARGETOS/$TARGETARCH" in \
      linux/amd64)  YTT_URL="https://github.com/carvel-dev/ytt/releases/download/${YTT_VERSION}/ytt-linux-amd64" ;; \
      linux/arm64)  YTT_URL="https://github.com/carvel-dev/ytt/releases/download/${YTT_VERSION}/ytt-linux-arm64" ;; \
      *) echo "Unsupported target: $TARGETOS/$TARGETARCH"; exit 1 ;; \
    esac; \
    apt-get update && apt-get install -y --no-install-recommends ca-certificates curl && rm -rf /var/lib/apt/lists/*; \
    curl -fsSL "$YTT_URL" -o /usr/local/bin/ytt; \
    chmod +x /usr/local/bin/ytt; \
    /usr/local/bin/ytt version

FROM $BASE_IMAGE
COPY --from=build /usr/local/bin/function /usr/local/bin/function
COPY --from=build /usr/local/bin/ytt /usr/local/bin/ytt
ENV PATH="/usr/local/bin:/usr/bin:/bin"
ENTRYPOINT ["function"]
