# download kubectl
FROM golang:1.23.3-alpine AS kubectl
RUN apk add --no-cache curl
RUN export VERSION=$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt) &&\
    export OS=$(go env GOOS) && \
    export ARCH=$(go env GOARCH) &&\
    curl -o /usr/local/bin/kubectl -L https://storage.googleapis.com/kubernetes-release/release/${VERSION}/bin/${OS}/${ARCH}/kubectl &&\
    chmod +x /usr/local/bin/kubectl

# build jsonnet-bundler
FROM golang:1.23.3-alpine AS jb
WORKDIR /tmp
RUN apk add --no-cache git make bash &&\
    git clone https://github.com/jsonnet-bundler/jsonnet-bundler &&\
    ls /bin &&\
    cd jsonnet-bundler &&\
    make static &&\
    mv _output/jb /usr/local/bin/jb

FROM golang:1.23.3-alpine AS helm
WORKDIR /tmp/helm
ARG HELM_VERSION
RUN apk add --no-cache jq curl
RUN export OS=$(go env GOOS) && \
    export ARCH=$(go env GOARCH) &&\
    if [[ -z ${HELM_VERSION} ]]; then export HELM_VERSION=$(curl --silent "https://api.github.com/repos/helm/helm/releases" | jq -r '.[0].tag_name'); fi && \
    curl -SL "https://get.helm.sh/helm-${HELM_VERSION}-${OS}-${ARCH}.tar.gz" > helm.tgz && \
    tar -xvf helm.tgz --strip-components=1

FROM golang:1.23.3-alpine AS kustomize
WORKDIR /tmp/kustomize
ARG KUSTOMIZE_VERSION
RUN apk add --no-cache jq curl
# Get the latest version of kustomize
# Releases are filtered by their name since the kustomize repository exposes multiple products in the releases
RUN export OS=$(go env GOOS) &&\
    export ARCH=$(go env GOARCH) &&\
    if [[ -z ${KUSTOMIZE_VERSION} ]]; then export KUSTOMIZE_VERSION=$(curl --silent "https://api.github.com/repos/kubernetes-sigs/kustomize/releases" | jq -r '[ .[] | select(.name | startswith("kustomize")) ] | .[0].tag_name | split("/")[1]'); fi && \
    echo "https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize/${KUSTOMIZE_VERSION}/kustomize_${KUSTOMIZE_VERSION}_${OS}_${ARCH}.tar.gz" && \
    curl -SL "https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize/${KUSTOMIZE_VERSION}/kustomize_${KUSTOMIZE_VERSION}_${OS}_${ARCH}.tar.gz" > kustomize.tgz && \
    tar -xvf kustomize.tgz

FROM golang:1.23.3 AS build
WORKDIR /app
COPY . .
RUN make static

# assemble final container
FROM alpine:3.20
RUN apk add --no-cache coreutils diffutils less git openssh-client && \
    apk upgrade --quiet
COPY --from=build /app/tk /usr/local/bin/tk
COPY --from=kubectl /usr/local/bin/kubectl /usr/local/bin/kubectl
COPY --from=jb /usr/local/bin/jb /usr/local/bin/jb
COPY --from=helm /tmp/helm/helm /usr/local/bin/helm
COPY --from=kustomize /tmp/kustomize/kustomize /usr/local/bin/kustomize
WORKDIR /app
ENTRYPOINT ["/usr/local/bin/tk"]
