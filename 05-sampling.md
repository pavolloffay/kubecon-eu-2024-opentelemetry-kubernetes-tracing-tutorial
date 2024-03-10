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

![OpenTelemetry Sampling](images/sampling-comparision.jpg)

### Head based sampling

Head sampling is a sampling technique used to make a sampling decision as early as possible. A decision to sample or drop a span or trace is not made by inspecting the trace as a whole.

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

Tail sampling is where the decision to sample a trace takes place by considering all or most of the spans within the trace. Tail Sampling gives you the option to sample your traces based on specific criteria derived from different parts of a trace, which isn’t an option with Head Sampling.

Usecase: Sample 100% of the traces that have an error-ing span in them.

```yaml
  processors: 
    tail_sampling:
      decision_wait: 10s # time to wait before making a sampling decision is made
      num_traces: 100 # number of traces to be kept in memory
      expected_new_traces_per_sec: 10 # expected rate of new traces per second
      policies:
        [          
          {
              name: keep-errors,
              type: status_code,
              status_code: {status_codes: [ERROR]}
            },
            {
              name: keep-slow-traces,
              type: latency,
              latency: {threshold_ms: 500}
            }
        ]
```

Applying this chart will start a new collector with the tailsampling processor

```shell
kubectl apply -f https://raw.githubusercontent.com/pavolloffay/kubecon-eu-2024-opentelemetry-kubernetes-tracing-tutorial/backend/03-tail-sampling-config.yaml
```

-----
### Advanced Topic: Sampling at scale with OpenTelemetry

Requires two deployments of the Collector, the first layer routing all the spans of a trace to the same collector in the downstream deployment (using load-balancing exporter). And the second layer doing the tail sampling.

![OpenTelemetry Sampling](images/scaling-otel-collector.jpg)

[excalidraw](https://excalidraw.com/#room=6a15d65ba4615c535a40,xcZD6DG977owHRoxpYY4Ag)

[Next steps](./06-RED-metrics.md)
