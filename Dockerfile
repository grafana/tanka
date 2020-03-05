FROM alpine

# less with `--RAW-CONTROL-CHARS` is required for tk show/ tk diff
RUN apk add --no-cache less

COPY tk /usr/local/bin/tk
ENTRYPOINT ["/usr/local/bin/tk"]


