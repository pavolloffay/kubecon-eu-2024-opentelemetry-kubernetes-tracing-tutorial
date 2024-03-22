# Exploring the Power of Distributed Tracing with OpenTelemetry on Kubernetes

This repository hosts content for tutorial for Kubecon EU 2024 Paris.

Previous tutorials:
* [Exploring the Power of OpenTelemetry on Kubernetes - Pavol Loffay, Benedikt Bongartz & Yuri Oliveira Sa, Red Hat; Severin Neumann, Cisco; Kristina Pathak, LightStep](https://github.com/pavolloffay/kubecon-eu-2023-opentelemetry-kubernetes-tutorial)
* [Tutorial: Exploring the Power of Metrics Collection with OpenTelemetry on Kubernetes - Pavol Loffay & Benedikt Bongartz, Red Hat; Anthony Mirabella, AWS; Matej Gera, Coralogix; Anusha Reddy Narapureddy, Apple](https://github.com/pavolloffay/kubecon-na-2023-opentelemetry-kubernetes-metrics-tutorial)

__Abstract__: Rolling out an observability solution is not a straightforward problem. There are many solutions and the final architecture can impact the effectiveness, robustness, and long-term maintenance aspects of the architecture. In this comprehensive tutorial, we will deploy an end-to-end distributed tracing stack on Kubernetes using the OpenTelemetry project. The tutorial will cover both manual and auto-instrumentation, extending the auto-instrumentation, collecting data with the OpenTelemetry collector and performing transformation on spans using OTTL, tail-based sampling, deriving metrics from traces, tracing with proxies/service meshes and collecting traces from Kubernetes infrastructure. After this session, the audience will be able to understand and use OpenTelemetry API/SDK, auto-instrumentation, collector, and operator to roll out a working distributed tracing stack on Kubernetes.

__Schedule__: https://sched.co/1YePA

__Slides__: [intro-slides](./intro-slides.pdf)

__Recording__: https://www.youtube.com/watch?v=nwy0I6vdtEE

## Agenda

Internal meeting doc: https://docs.google.com/document/d/1rbc0JqMP7i4koKpxqb9gYovmAlJ_BRN1Ttg3EhY9cbY/edit

Each tutorial step is located in a separate file:

1. [Welcome & Setup](01-welcome-setup.md) (Pavol, 5 min)
1. [OpenTelemetry distributed tracing introduction](02-tracing-introduction.md) (Matej, 10 min)
1. [Auto-instrumentation](03-auto-instrumentation.md) (Pavol, 25 min)
1. [Manual-instrumentation](04-manual-instrumentation.md) (Bene, 10 min)
1. [Sampling](05-sampling.md) (Bene & Anu, 15 min)
1. [Metrics from Traces](06-RED-metrics.md) (Anthony, 10 min)
1. [OpenTelemetry Transformation Language and Spans](07-ottl.md) (Matej, 10 min)
1. Wrap up & Questions
1. [K8S-Tracing](08-k8s-tracing.md) (Bene, optional)
