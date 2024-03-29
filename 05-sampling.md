# Sampling

This tutorial step covers the basic usage of the OpenTelemetry Collector on Kubernetes and how to reduce costs using sampling techniques.

## Overview

In chapter 3 we saw the [schematic structure of the dice game application](https://github.com/pavolloffay/kubecon-eu-2024-opentelemetry-kubernetes-tracing-tutorial/blob/main/03-auto-instrumentation.md#application-description). The following diagram illustrates how the telemetry data collected there is exported and stored. [excalidraw](https://excalidraw.com/#json=15BrdSOMEkc9RA5cxeqwz,urTmfk01mbx7V-PpQI7KgA)

![tracing setup](images/tracing-setup.png)

## Sampling, what does it mean and why is it important?

Sampling refers to the practice of selectively capturing and recording traces of requests flowing through a distributed system, rather than capturing every single request. It is crucial in distributed tracing systems because modern distributed applications often generate a massive volume of requests and transactions, which can overwhelm the tracing infrastructure or lead to excessive storage costs if every request is traced in detail.

For example, a medium sized setup producing ~1M traces per minute can result in a cost of approximately $250,000 per month. (Note that this depends on your infrastructure costs, the SaaS provider you choose, the amount of metadata, etc.) You may want to check some service costs to get a better idea.

Pricing:
- AWS Xray ([calculator](https://aws.amazon.com/xray/pricing/))
- GCP Cloud Trace ([pricing](https://cloud.google.com/stackdriver/pricing#trace-costs))

```
GCP

Feature           Price                 Free allotment per month  Effective date
Trace ingestion   $0.20/million spans   First 2.5 million spans   November 1, 2018 
---

X-Ray Tracing

Traces recorded cost $5.00 per 1 million traces recorded ($0.000005 per trace).

Traces retrieved cost $0.50 per 1 million traces retrieved ($0.0000005 per trace).

Traces scanned cost $0.50 per 1 million traces scanned ($0.0000005 per trace).

X-Ray Insights traces stored costs $1.00 per million traces recorded ($0.000001 per trace).
```

For more details, check the [offical documentation](https://opentelemetry.io/docs/concepts/sampling/).

### How can we now reduce the number of traces?

![OpenTelemetry Sampling](images/sampling-venn.svg)

### Comparing Sampling Approaches

![OpenTelemetry Sampling](images/sampling-comparision.jpg)

### How to implement head sampling with OpenTelemetry

Head sampling is a sampling technique used to make a sampling decision as early as possible. A decision to sample or drop a span or trace is not made by inspecting the trace as a whole.

For the list of all available samplers, check the [offical documentation](https://opentelemetry.io/docs/languages/sdk-configuration/general/#otel_traces_sampler)

#### Auto Instrumentation

Update the sampling % in the Auto Instrumentation CR and restart the deployment for the configurations to take effect.

https://github.com/pavolloffay/kubecon-eu-2024-opentelemetry-kubernetes-tracing-tutorial/blob/d4b917c1cc4a411f59ae5dd770b22de1de9f6020/app/instrumentation-head-sampling.yaml#L13-L15

```yaml
kubectl apply -f https://raw.githubusercontent.com/pavolloffay/kubecon-eu-2024-opentelemetry-kubernetes-tracing-tutorial/main/app/instrumentation-head-sampling.yaml
kubectl rollout restart deployment.apps/backend1-deployment -n tutorial-application
kubectl get pods -w -n tutorial-application
```

Describe the pod spec for the backend1 deployment to see the updated sampling rate.

```bash
kubectl describe pod backend1-deployment-64ddcc76fd-w85zh -n tutorial-application
```

```diff
    Environment:
          OTEL_TRACES_SAMPLER:                 parentbased_traceidratio
-         OTEL_TRACES_SAMPLER_ARG:             1
+         OTEL_TRACES_SAMPLER_ARG:             0.5
```

This tells the SDK to sample spans such that only 50% of traces get created.

#### Manual Instrumentation 

You can also configure the ParentBasedTraceIdRatioSampler in code.A [`Sampler`](https://pkg.go.dev/go.opentelemetry.io/otel/sdk/trace#Sampler) can be set on the tracer provider using the [`WithSampler`](https://pkg.go.dev/go.opentelemetry.io/otel/sdk/trace#WithSampler)
option, as follows:

```go
provider := trace.NewTracerProvider(
    trace.WithSampler(trace.NewParentBasedTraceIdRatioSampler(0.5)),
)
```

### How to implement tail sampling in the OpenTelemetry Collector

Tail sampling is where the decision to sample a trace takes place by considering all or most of the spans within the trace. Tail Sampling gives you the option to sample your traces based on specific criteria derived from different parts of a trace, which isn’t an option with Head Sampling.

Update the ENV variables below in the backend2 deployment, which generates random spans with errors and high latencies.

```shell
kubectl set env deployment backend2-deployment RATE_ERROR=50 RATE_HIGH_DELAY=50 -n tutorial-application 
kubectl get pods -n tutorial-application -w
```

Deploy the opentelemetry collector with `tail_sampling` enabled.

```shell
kubectl apply -f https://raw.githubusercontent.com/pavolloffay/kubecon-eu-2024-opentelemetry-kubernetes-tracing-tutorial/main/backend/05-collector-1.yaml
kubectl get pods -n observability-backend -w
```

Now, let’s walk-through the tail-sampling processor configuration, placed in the `processors` section of the collector configuration file:

```yaml
  # 1. Sample 100% of traces with ERROR-ing spans
  # 2. Sample 100% of trace which have a duration longer than 500ms
  # 3. Randomized sampling of 10% of traces without errors and latencies.
  processors: 
    tail_sampling:
      decision_wait: 10s # time to wait before making a sampling decision
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
            },
            {
              name: randomized-policy,
              type: probabilistic,
              probabilistic: {sampling_percentage: 10}
            }
        ]
```

Now let's execute some requests on the app [http://localhost:4000/](http://localhost:4000/) and see traces in the Jaeger console [http://localhost:16686/](http://localhost:16686/).

The image next is an example of what you might see in your backend with this sample configuration. With this configuration, you’ll get all traces with errors and latencies exceeding 500ms, as well as a random sample of other traces based on the rate we’ve configured.

![OpenTelemetry Sampling](images/jaeger-tail-sampling.jpg)

You also have the flexibility to add other policies. For the list of all policies, check the [offical documentation](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/processor/tailsamplingprocessor/README.md)

Here are a few examples:

- `always_sample`: Sample all traces.
- `string_attribute`: Sample based on string attribute values, both exact and
  regular expression value matches are supported. For example, you could sample
  based on specific custom attribute values.

-----
### Advanced Topic: Tail Sampling at scale with OpenTelemetry
> [!NOTE]  
> This is an optional more advanced section.

All spans of a trace must be processed by the same collector for tail sampling to function properly, posing scalability challenges. Initially, a single collector may suffice, but as the system grows, a two-layer setup becomes necessary. It requires two deployments of the collector, with the first layer routing all spans of a trace to the same collector in the downstream deployment (using a [load-balancing exporter](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/exporter/loadbalancingexporter/README.md)), and the second layer performing the tail sampling.

![OpenTelemetry Sampling](images/scaling-otel-collector.jpg)

[excalidraw](https://excalidraw.com/#room=6a15d65ba4615c535a40,xcZD6DG977owHRoxpYY4Ag)

Apply the YAML below to deploy a layer of Collectors containing the load-balancing exporter in front of collectors performing tail-sampling:

```shell
kubectl apply -f https://raw.githubusercontent.com/pavolloffay/kubecon-eu-2024-opentelemetry-kubernetes-tracing-tutorial/main/backend/05-collector-2.yaml
kubectl get pods -n observability-backend -w
```

```bash
jaeger-bc5f49d78-627ct                    1/1     Running   0          100m
otel-collector-b48b5d66d-k5dsc            1/1     Running   0          4m42s
otel-gateway-collector-0                  1/1     Running   0          3m38s
otel-gateway-collector-1                  1/1     Running   0          3m38s
prometheus-77f88ccf7f-dfwh2               1/1     Running   0          100m

```

Now, let’s walk-through the load-balancing exporter configuration, placed in the `exporters` section of the collector (layer 1) configuration file:

```yaml
  exporters:
    debug:
    # routing_key property is used to route spans to exporters based on traceID/service name
    loadbalancing:
      routing_key: "traceID"
      protocol:
        otlp:
          timeout: 1s
          tls:
            insecure: true
      resolver:
        k8s:
          service: otel-gateway.observability-backend
          ports: 
            - 4317
```

### Advanced Topic: Jaeger's Remote Sampling extension
> [!NOTE]  
> This is an optional more advanced section.
 
This extension allows serving sampling strategies following the Jaeger's remote sampling API. This extension can be configured to proxy requests to a backing remote sampling server, which could potentially be a Jaeger Collector down the pipeline, or a static JSON file from the local file system.

#### Example Configuration

```yaml
extensions:
  jaegerremotesampling:
    source:
      reload_interval: 30s
      remote:
        endpoint: jaeger-collector:14250
  jaegerremotesampling/1:
    source:
      reload_interval: 1s
      file: /etc/otelcol/sampling_strategies.json
  jaegerremotesampling/2:
    source:
      reload_interval: 1s
      file: http://jaeger.example.com/sampling_strategies.json
```

For more details, check the [offical documentation](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/extension/jaegerremotesampling/README.md)


[Next steps](./06-RED-metrics.md)
