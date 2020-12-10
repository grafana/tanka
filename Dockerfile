# download kubectl
FROM alpine as kubectl
RUN apk add --no-cache curl
RUN export VERSION=$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt) &&\
    curl -o /usr/local/bin/kubectl -L https://storage.googleapis.com/kubernetes-release/release/${VERSION}/bin/linux/amd64/kubectl &&\
    chmod +x /usr/local/bin/kubectl

# build jsonnet-bundler
FROM golang as jb
WORKDIR /tmp
RUN git clone https://github.com/jsonnet-bundler/jsonnet-bundler &&\
    cd jsonnet-bundler &&\
    make static &&\
    mv _output/jb /usr/local/bin/jb

FROM alpine as helm
RUN apk add --no-cache curl bash openssl
RUN curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 &&\
    chmod 700 get_helm.sh &&\
    ./get_helm.sh

# assemble final container
FROM alpine
RUN apk add --no-cache coreutils diffutils less git openssh-client
COPY tk /usr/local/bin/tk
COPY --from=kubectl /usr/local/bin/kubectl /usr/local/bin/kubectl
COPY --from=jb /usr/local/bin/jb /usr/local/bin/jb
COPY --from=helm /usr/local/bin/helm /usr/local/bin/helm
WORKDIR /app
ENTRYPOINT ["/usr/local/bin/tk"]
