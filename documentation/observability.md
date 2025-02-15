# Observability

Scheduler currently supports logging and metrics. Both are exported via the OpenTelemetry protocol (GRPC) and can be
collected by any OpenTelemetry-compatible collector.

## Logging

Logging can be configured via the `LOG_LEVEL` environment variable with `info`,`debug`, `warn` or `error` levels.

## Metrics

Metrics can be enabled by setting the `METRICS_ENABLED` environment variable to `true`. Metrics are exported via the
OpenTelemetry protocol (GRPC).

The following manager metrics are currently exported:

- `http_requests_total`: The total number of HTTP requests received by the server.
- `http_request_duration_seconds`: The duration of HTTP requests in seconds.
- `http_errors_total`: The total number of failed HTTP requests.

The following runner metrics are currently exported:

- `scheduler_jobs_total`: The total number of jobs that have been scheduled.
- `scheduler_jobs_failed_total`: The total number of jobs that have failed.
- `scheduler_jobs_duration_seconds`: The duration of jobs in seconds.
- `scheduler_jobs_in_execution`: The total number of jobs currently in execution.
