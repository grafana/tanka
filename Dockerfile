# download kubectl
FROM golang:1.26.4-alpine@sha256:3ad57304ad93bbec8548a0437ad9e06a455660655d9af011d58b993f6f615648 AS kubectl
ARG KUBECTL_VERSION=1.34.9
RUN apk add --no-cache curl
RUN export OS=$(go env GOOS) && \
    export ARCH=$(go env GOARCH) &&\
    curl -o /usr/local/bin/kubectl -L https://dl.k8s.io/release/v${KUBECTL_VERSION}/bin/${OS}/${ARCH}/kubectl &&\
    chmod +x /usr/local/bin/kubectl

# build jsonnet-bundler
FROM golang:1.26.4-alpine@sha256:3ad57304ad93bbec8548a0437ad9e06a455660655d9af011d58b993f6f615648 AS jb
WORKDIR /tmp
RUN apk add --no-cache git make bash &&\
    git clone https://github.com/jsonnet-bundler/jsonnet-bundler &&\
    ls /bin &&\
    cd jsonnet-bundler &&\
    make static &&\
    mv _output/jb /usr/local/bin/jb

FROM golang:1.26.4-alpine@sha256:3ad57304ad93bbec8548a0437ad9e06a455660655d9af011d58b993f6f615648 AS helm
WORKDIR /tmp/helm
ARG HELM_VERSION=3.21.1
RUN apk add --no-cache jq curl
RUN export OS=$(go env GOOS) && \
    export ARCH=$(go env GOARCH) &&\
    curl -SL "https://get.helm.sh/helm-v${HELM_VERSION}-${OS}-${ARCH}.tar.gz" > helm.tgz && \
    tar -xvf helm.tgz --strip-components=1

FROM golang:1.26.4-alpine@sha256:3ad57304ad93bbec8548a0437ad9e06a455660655d9af011d58b993f6f615648 AS kustomize
WORKDIR /tmp/kustomize
ARG KUSTOMIZE_VERSION=5.8.1
RUN apk add --no-cache jq curl
RUN export OS=$(go env GOOS) &&\
    export ARCH=$(go env GOARCH) &&\
    echo "https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize/v${KUSTOMIZE_VERSION}/kustomize_v${KUSTOMIZE_VERSION}_${OS}_${ARCH}.tar.gz" && \
    curl -SL "https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize/v${KUSTOMIZE_VERSION}/kustomize_v${KUSTOMIZE_VERSION}_${OS}_${ARCH}.tar.gz" > kustomize.tgz && \
    tar -xvf kustomize.tgz

FROM golang:1.26.4@sha256:32c0e6e5c4f6707717051091b4d0b077464a679eaab563e11474efc5328e2aa5 AS build
WORKDIR /app
COPY . .
RUN make static

# assemble final container
FROM alpine:3.24.0@sha256:a2d49ea686c2adfe3c992e47dc3b5e7fa6e6b5055609400dc2acaeb241c829f4
RUN apk add --no-cache coreutils diffutils less git openssh-client && \
    apk upgrade --quiet
COPY --from=build /app/tk /usr/local/bin/tk
COPY --from=kubectl /usr/local/bin/kubectl /usr/local/bin/kubectl
COPY --from=jb /usr/local/bin/jb /usr/local/bin/jb
COPY --from=helm /tmp/helm/helm /usr/local/bin/helm
COPY --from=kustomize /tmp/kustomize/kustomize /usr/local/bin/kustomize
WORKDIR /app
ENTRYPOINT ["/usr/local/bin/tk"]
