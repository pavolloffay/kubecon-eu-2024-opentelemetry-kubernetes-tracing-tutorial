#!/bin/sh

if [[ "${OTEL_INSTRUMENTATION_ENABLED}" == "true" ]] ; then
    echo 'Run with instrumentation'
    env OTEL_SERVICE_NAME=${OTEL_SERVICE_NAME:-backend2} \
    OTEL_TRACES_EXPORTER=${OTEL_TRACES_EXPORTER:-logging} \
    OTEL_METRICS_EXPORTER=${OTEL_METRICS_EXPORTER:-logging} \
    OTEL_LOGS_EXPORTER=${OTEL_LOGS_EXPORTER:-logging} \
    java -javaagent:./javaagent.jar -jar ./build/libs/dice-0.0.1-SNAPSHOT.jar
else 
    java -jar ./build/libs/dice-0.0.1-SNAPSHOT.jar
fi
