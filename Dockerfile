# download kubectl
FROM golang:1.25.1-alpine@sha256:b6ed3fd0452c0e9bcdef5597f29cc1418f61672e9d3a2f55bf02e7222c014abd AS kubectl
ARG KUBECTL_VERSION=1.34.1
RUN apk add --no-cache curl
RUN export OS=$(go env GOOS) && \
    export ARCH=$(go env GOARCH) &&\
    curl -o /usr/local/bin/kubectl -L https://cdn.dl.k8s.io/release/v${KUBECTL_VERSION}/bin/${OS}/${ARCH}/kubectl &&\
    chmod +x /usr/local/bin/kubectl

# build jsonnet-bundler
FROM golang:1.25.1-alpine@sha256:b6ed3fd0452c0e9bcdef5597f29cc1418f61672e9d3a2f55bf02e7222c014abd AS jb
WORKDIR /tmp
RUN apk add --no-cache git make bash &&\
    git clone https://github.com/jsonnet-bundler/jsonnet-bundler &&\
    ls /bin &&\
    cd jsonnet-bundler &&\
    make static &&\
    mv _output/jb /usr/local/bin/jb

FROM golang:1.25.1-alpine@sha256:b6ed3fd0452c0e9bcdef5597f29cc1418f61672e9d3a2f55bf02e7222c014abd AS helm
WORKDIR /tmp/helm
ARG HELM_VERSION=3.19.0
RUN apk add --no-cache jq curl
RUN export OS=$(go env GOOS) && \
    export ARCH=$(go env GOARCH) &&\
    curl -SL "https://get.helm.sh/helm-v${HELM_VERSION}-${OS}-${ARCH}.tar.gz" > helm.tgz && \
    tar -xvf helm.tgz --strip-components=1

FROM golang:1.25.1-alpine@sha256:b6ed3fd0452c0e9bcdef5597f29cc1418f61672e9d3a2f55bf02e7222c014abd AS kustomize
WORKDIR /tmp/kustomize
ARG KUSTOMIZE_VERSION=5.7.1
RUN apk add --no-cache jq curl
RUN export OS=$(go env GOOS) &&\
    export ARCH=$(go env GOARCH) &&\
    echo "https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize/v${KUSTOMIZE_VERSION}/kustomize_v${KUSTOMIZE_VERSION}_${OS}_${ARCH}.tar.gz" && \
    curl -SL "https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize/v${KUSTOMIZE_VERSION}/kustomize_v${KUSTOMIZE_VERSION}_${OS}_${ARCH}.tar.gz" > kustomize.tgz && \
    tar -xvf kustomize.tgz

FROM golang:1.25.1@sha256:bb979b278ffb8d31c8b07336fd187ef8fafc8766ebeaece524304483ea137e96 AS build
WORKDIR /app
COPY . .
RUN make static

# assemble final container
FROM alpine:3.22@sha256:4bcff63911fcb4448bd4fdacec207030997caf25e9bea4045fa6c8c44de311d1
RUN apk add --no-cache coreutils diffutils less git openssh-client && \
    apk upgrade --quiet
COPY --from=build /app/tk /usr/local/bin/tk
COPY --from=kubectl /usr/local/bin/kubectl /usr/local/bin/kubectl
COPY --from=jb /usr/local/bin/jb /usr/local/bin/jb
COPY --from=helm /tmp/helm/helm /usr/local/bin/helm
COPY --from=kustomize /tmp/kustomize/kustomize /usr/local/bin/kustomize
WORKDIR /app
ENTRYPOINT ["/usr/local/bin/tk"]
