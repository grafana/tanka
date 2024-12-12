# base stage for downloading binaries
FROM golang:1.23.4-alpine3.21 AS base
RUN apk add --no-cache curl


# download kubectl
FROM base AS kubectl
# renovate: datasource=github-releases packageName=kubernetes/kubernetes
ARG KUBECTL_VERSION=1.32.0
RUN export OS=$(go env GOOS) && \
    export ARCH=$(go env GOARCH) &&\
    curl -o /usr/local/bin/kubectl -L https://cdn.dl.k8s.io/release/v${KUBECTL_VERSION}/bin/${OS}/${ARCH}/kubectl &&\
    chmod +x /usr/local/bin/kubectl

# download jsonnet-bundler
FROM base AS jb
WORKDIR /tmp
# renovate: datasource=github-releases packageName=jsonnet-bundler/jsonnet-bundler
ARG JB_VERSION=0.5.1
RUN apk add --no-cache curl && \
    OS=$(go env GOOS) && \
    ARCH=$(go env GOARCH) && \
    curl -o /usr/local/bin/jb -SL "https://github.com/jsonnet-bundler/jsonnet-bundler/releases/download/v${JB_VERSION}/jb-${OS}-${ARCH}" && \
    chmod +x /usr/local/bin/jb

# download helm
FROM base AS helm
WORKDIR /tmp/helm
# renovate: datasource=github-releases packageName=helm/helm
ARG HELM_VERSION=3.16.3
RUN apk add --no-cache curl
RUN export OS=$(go env GOOS) && \
    export ARCH=$(go env GOARCH) &&\
    curl -SL "https://get.helm.sh/helm-v${HELM_VERSION}-${OS}-${ARCH}.tar.gz" > helm.tgz && \
    tar -xvf helm.tgz --strip-components=1

# download kustomize
FROM base AS kustomize
WORKDIR /tmp/kustomize
# renovate: datasource=github-releases packageName=kubernetes-sigs/kustomize
ARG KUSTOMIZE_VERSION=5.5.0
RUN apk add --no-cache curl
RUN export OS=$(go env GOOS) &&\
    export ARCH=$(go env GOARCH) &&\
    echo "https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize/v${KUSTOMIZE_VERSION}/kustomize_v${KUSTOMIZE_VERSION}_${OS}_${ARCH}.tar.gz" && \
    curl -SL "https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize/v${KUSTOMIZE_VERSION}/kustomize_v${KUSTOMIZE_VERSION}_${OS}_${ARCH}.tar.gz" > kustomize.tgz && \
    tar -xvf kustomize.tgz

FROM base AS build
WORKDIR /app
COPY . .
RUN make static

# assemble final container
FROM alpine:3.21.0
RUN apk add --no-cache coreutils diffutils less git openssh-client && \
    apk upgrade --quiet
COPY --from=build /app/tk /usr/local/bin/tk
COPY --from=kubectl /usr/local/bin/kubectl /usr/local/bin/kubectl
COPY --from=jb /usr/local/bin/jb /usr/local/bin/jb
COPY --from=helm /tmp/helm/helm /usr/local/bin/helm
COPY --from=kustomize /tmp/kustomize/kustomize /usr/local/bin/kustomize
WORKDIR /app
ENTRYPOINT ["/usr/local/bin/tk"]
