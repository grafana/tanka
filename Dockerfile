# download kubectl
FROM golang:alpine3.13 as kubectl
RUN apk add --no-cache curl
RUN export VERSION=$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt) &&\
    export OS=$(go env GOOS) && \
    export ARCH=$(go env GOARCH) &&\
    curl -o /usr/local/bin/kubectl -L https://storage.googleapis.com/kubernetes-release/release/${VERSION}/bin/${OS}/${ARCH}/kubectl &&\
    chmod +x /usr/local/bin/kubectl

# build jsonnet-bundler
FROM golang:alpine3.13 as jb
WORKDIR /tmp
RUN apk add --no-cache git make bash &&\
    git clone https://github.com/jsonnet-bundler/jsonnet-bundler &&\
    ls /bin &&\
    cd jsonnet-bundler &&\
    make static &&\
    mv _output/jb /usr/local/bin/jb

FROM golang:alpine3.13 as helm
WORKDIR /tmp/helm
RUN apk add --no-cache jq curl
RUN export TAG=$(curl --silent "https://api.github.com/repos/helm/helm/releases/latest" | jq -r .tag_name) &&\
    export OS=$(go env GOOS) &&\
    export ARCH=$(go env GOARCH) &&\
    curl -SL "https://get.helm.sh/helm-${TAG}-${OS}-${ARCH}.tar.gz" > helm.tgz && \
    tar -xvf helm.tgz --strip-components=1

FROM golang:alpine3.13 as kustomize
WORKDIR /tmp/kustomize
RUN apk add --no-cache jq curl
# Get the latest version of kustomize
# Releases are filtered by their name since the kustomize repository exposes multiple products in the releases
RUN export TAG=$(curl --silent "https://api.github.com/repos/kubernetes-sigs/kustomize/releases" | jq -r '[ .[] | select(.name | startswith("kustomize")) ] | .[0].tag_name') &&\
    export VERSION_TAG=${TAG#*/} &&\
    export OS=$(go env GOOS) &&\
    export ARCH=$(go env GOARCH) &&\
    echo "https://github.com/kubernetes-sigs/kustomize/releases/download/${TAG}/kustomize_${VERSION_TAG}_${OS}_${ARCH}.tar.gz" && \
    curl -SL "https://github.com/kubernetes-sigs/kustomize/releases/download/${TAG}/kustomize_${VERSION_TAG}_${OS}_${ARCH}.tar.gz" > kustomize.tgz && \
    tar -xvf kustomize.tgz

# assemble final container
FROM alpine:3.13
RUN apk add --no-cache coreutils diffutils less git openssh-client
COPY tk /usr/local/bin/tk
COPY --from=kubectl /usr/local/bin/kubectl /usr/local/bin/kubectl
COPY --from=jb /usr/local/bin/jb /usr/local/bin/jb
COPY --from=helm /tmp/helm/helm /usr/local/bin/helm
COPY --from=kustomize /tmp/kustomize/kustomize /usr/local/bin/kustomize
WORKDIR /app
ENTRYPOINT ["/usr/local/bin/tk"]
