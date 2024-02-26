#!/bin/sh

if [[ "${OTEL_INSTRUMENTATION_ENABLED}" == "true" ]] ; then
    echo 'Run with instrumentation'
    env OTEL_SERVICE_NAME=${OTEL_SERVICE_NAME:-frontend} \
    OTEL_TRACES_EXPORTER=${OTEL_TRACES_EXPORTER:-console} \
    OTEL_METRICS_EXPORTER=${OTEL_METRICS_EXPORTER:-console} \
    OTEL_LOGS_EXPORTER=${OTEL_LOGS_EXPORTER:-console} \
    node -r ./instrument.js index.js
else 
    node index.js
fi