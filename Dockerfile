# download kubectl
FROM alpine as kubectl
RUN apk add --no-cache curl
RUN export VERSION=$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt) &&\
    curl -o /usr/local/bin/kubectl -L https://storage.googleapis.com/kubernetes-release/release/${VERSION}/bin/linux/amd64/kubectl &&\
    chmod +x /usr/local/bin/kubectl

# build jsonnet-bundler
FROM golang:alpine as jb
WORKDIR /tmp
RUN apk add --no-cache git make bash &&\
    git clone https://github.com/jsonnet-bundler/jsonnet-bundler &&\
    cd jsonnet-bundler &&\
    make static &&\
    mv _output/jb /usr/local/bin/jb

FROM golang:alpine as helm
WORKDIR /tmp/helm
RUN apk add --no-cache jq curl
RUN export TAG=$(curl --silent "https://api.github.com/repos/helm/helm/releases/latest" | jq -r .tag_name) &&\
    export OS=$(go env GOOS) &&\
    export ARCH=$(go env GOARCH) &&\
    curl -SL "https://get.helm.sh/helm-$TAG-$OS-$ARCH.tar.gz" > helm.tgz && \
    tar -xvf helm.tgz --strip-components=1

# assemble final container
FROM alpine
RUN apk add --no-cache coreutils diffutils less git openssh-client
COPY tk /usr/local/bin/tk
COPY --from=kubectl /usr/local/bin/kubectl /usr/local/bin/kubectl
COPY --from=jb /usr/local/bin/jb /usr/local/bin/jb
COPY --from=helm /tmp/helm/helm /usr/local/bin/helm
WORKDIR /app
ENTRYPOINT ["/usr/local/bin/tk"]
