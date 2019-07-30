FROM alpine
COPY tk /usr/local/bin/tk
ENTRYPOINT ["/usr/local/bin/tk"]


