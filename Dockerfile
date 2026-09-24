# download kubectl
FROM golang:1.27.1-alpine@sha256:8a5910f31396cd4d89662f56c68b3ae31d374308270a1c3bd96672ee5ed43414 AS kubectl
ARG KUBECTL_VERSION=1.34.11
RUN apk add --no-cache curl
RUN export OS=$(go env GOOS) && \
    export ARCH=$(go env GOARCH) &&\
    curl -o /usr/local/bin/kubectl -L https://dl.k8s.io/release/v${KUBECTL_VERSION}/bin/${OS}/${ARCH}/kubectl &&\
    chmod +x /usr/local/bin/kubectl

# build jsonnet-bundler
FROM golang:1.27.1-alpine@sha256:8a5910f31396cd4d89662f56c68b3ae31d374308270a1c3bd96672ee5ed43414 AS jb
WORKDIR /tmp
RUN apk add --no-cache git make bash &&\
    git clone https://github.com/jsonnet-bundler/jsonnet-bundler &&\
    ls /bin &&\
    cd jsonnet-bundler &&\
    make static &&\
    mv _output/jb /usr/local/bin/jb

FROM golang:1.27.1-alpine@sha256:8a5910f31396cd4d89662f56c68b3ae31d374308270a1c3bd96672ee5ed43414 AS helm
WORKDIR /tmp/helm
ARG HELM_VERSION=3.21.1
RUN apk add --no-cache jq curl
RUN export OS=$(go env GOOS) && \
    export ARCH=$(go env GOARCH) &&\
    curl -SL "https://get.helm.sh/helm-v${HELM_VERSION}-${OS}-${ARCH}.tar.gz" > helm.tgz && \
    tar -xvf helm.tgz --strip-components=1

FROM golang:1.27.1-alpine@sha256:8a5910f31396cd4d89662f56c68b3ae31d374308270a1c3bd96672ee5ed43414 AS kustomize
WORKDIR /tmp/kustomize
ARG KUSTOMIZE_VERSION=5.8.1
RUN apk add --no-cache jq curl
RUN export OS=$(go env GOOS) &&\
    export ARCH=$(go env GOARCH) &&\
    echo "https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize/v${KUSTOMIZE_VERSION}/kustomize_v${KUSTOMIZE_VERSION}_${OS}_${ARCH}.tar.gz" && \
    curl -SL "https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize/v${KUSTOMIZE_VERSION}/kustomize_v${KUSTOMIZE_VERSION}_${OS}_${ARCH}.tar.gz" > kustomize.tgz && \
    tar -xvf kustomize.tgz

FROM golang:1.27.1@sha256:f44f6e88636cfb311f9ebace870ded69d943f227bb3cb27d32ffd84ea18c43ea AS build
WORKDIR /app
COPY . .
RUN make static

# assemble final container
FROM alpine:3.24.2@sha256:294b683cb724975bec92580e1e685676bd4b50bda910ddb8c51d4cabeaec77e6
RUN apk add --no-cache coreutils diffutils less git openssh-client && \
    apk upgrade --quiet
COPY --from=build /app/tk /usr/local/bin/tk
COPY --from=kubectl /usr/local/bin/kubectl /usr/local/bin/kubectl
COPY --from=jb /usr/local/bin/jb /usr/local/bin/jb
COPY --from=helm /tmp/helm/helm /usr/local/bin/helm
COPY --from=kustomize /tmp/kustomize/kustomize /usr/local/bin/kustomize
WORKDIR /app
ENTRYPOINT ["/usr/local/bin/tk"]
