#!/bin/bash

PORT=${PORT:-5000}
HOST=${HOST:-"0.0.0.0"}

if [[ "${OTEL_INSTRUMENTATION_ENABLED}" == "true" ]] ; then
    echo 'Run with instrumentation'
    env OTEL_SERVICE_NAME=${OTEL_SERVICE_NAME:-backend1} \
    OTEL_TRACES_EXPORTER=${OTEL_TRACES_EXPORTER:-console} \
    OTEL_METRICS_EXPORTER=${OTEL_METRICS_EXPORTER:-console} \
    OTEL_LOGS_EXPORTER=${OTEL_LOGS_EXPORTER:-console} \
    OTEL_PYTHON_LOG_CORRELATION=${OTEL_PYTHON_LOG_CORRELATION:-"true"} \
    opentelemetry-instrument python3 -m flask run --host="${HOST}" --port="${PORT}"
else 
    python3 -m flask run --host="${HOST}" --port "${PORT}"
fi