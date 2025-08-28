---
title: Telemetry
---

Tanka supports sending OpenTelemetry traces to a `http/protobuf` endpoint.
To use this, export the following environment variable:

```
export OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4318
```

Note that we try to keep traces to the critical paths around the `export` command since these are usually the areas where performance is most important in automated workflows.
