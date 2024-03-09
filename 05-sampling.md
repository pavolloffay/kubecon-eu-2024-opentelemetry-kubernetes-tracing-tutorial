# Sampling

This tutorial step covers the basic usage of the OpenTelemetry Collector on Kubernetes and how to reduce costs using sampling techniques.

## Overview

![tracing setup](images/tracing-setup.png)

[excalidraw](https://excalidraw.com/#json=15BrdSOMEkc9RA5cxeqwz,urTmfk01mbx7V-PpQI7KgA)

### OpenTelemetry Collector on k8s

After installing the OpenTelemetry Operator, the `v1alpha1.OpenTelemetryCollector` simplifies the operation of the OpenTelemetry Collector on Kubernetes. There are different deployment modes available, breaking config changes are migrated automatically, provides integration with Prometheus (including operating on Prometheus Operator CRs) and simplifies sidecar injection.

TODO: update collector
```yaml

```

## Sampling, what does it mean and why is it important?

Sampling refers to the practice of selectively capturing and recording traces of requests flowing through a distributed system, rather than capturing every single request. It is crucial in distributed tracing systems because modern distributed applications often generate a massive volume of requests and transactions, which can overwhelm the tracing infrastructure or lead to excessive storage costs if every request is
traced in detail.

For example, a medium sized setup producing ~1M traces per minute can result in a cost of approximately $250,000 per month. (Note that this depends on your infrastructure costs, the SaaS provider you choose, the amount of metadata, etc.)

To get a better feel for the cost, you may want to play with some SaaS cost calculators.

- TODO
- TODO
- TODO

For more details, check the [offical documentation](https://opentelemetry.io/docs/concepts/sampling/).

### How can we now reduce the number of traces?

![OpenTelemetry Sampling](images/sampling-venn.svg)

### Comparing Sampling Approaches

#### Head based sampling

TODO: Update SDK config
```yaml

```

<details>
  <summary>Jaeger's Remote Sampling extension</summary>
 
TODO:
https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/extension/jaegerremotesampling/README.md


https://opentelemetry.io/docs/languages/sdk-configuration/general/#otel_traces_sampler_arg
 
</details>

### Tailbased Sampling

TODO

---

[Next steps](./06-RED-metrics.md)


