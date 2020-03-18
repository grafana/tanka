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

# assemble final container
FROM alpine
RUN apk add --no-cache coreutils diffutils less git
COPY tk /usr/local/bin/tk
COPY --from=kubectl /usr/local/bin/kubectl /usr/local/bin/kubectl
COPY --from=jb /usr/local/bin/jb /usr/local/bin/jb
WORKDIR /app
ENTRYPOINT ["/usr/local/bin/tk"]
