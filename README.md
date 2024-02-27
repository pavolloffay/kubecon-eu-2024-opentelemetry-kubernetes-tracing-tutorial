# Exploring the Power of Distributed Tracing with OpenTelemetry on Kubernetes

This repository hosts content for tutorial for Kubecon EU 2024 Paris.

__Abstract__: Rolling out an observability solution is not a straightforward problem. There are many solutions and the final architecture can impact the effectiveness, robustness, and long-term maintenance aspects of the architecture. In this comprehensive tutorial, we will deploy an end-to-end distributed tracing stack on Kubernetes using the OpenTelemetry project. The tutorial will cover both manual and auto-instrumentation, extending the auto-instrumentation, collecting data with the OpenTelemetry collector and performing transformation on spans using OTTL, tail-based sampling, deriving metrics from traces, tracing with proxies/service meshes and collecting traces from Kubernetes infrastructure. After this session, the audience will be able to understand and use OpenTelemetry API/SDK, auto-instrumentation, collector, and operator to roll out a working distributed tracing stack on Kubernetes.

__Schedule__: https://kccnceu2024.sched.com/event/1YePAr

__Slides__: https://docs.google.com/presentation/d/1IdQgrDKp4cP6fJ_IHlv1QszRMGL_uquW/edit#slide=id.p2

__Recording__: 

## Agenda

Internal meeting doc: https://docs.google.com/document/d/1rbc0JqMP7i4koKpxqb9gYovmAlJ_BRN1Ttg3EhY9cbY/edit

Each tutorial step is located in a separate file:

1. [Welcome & Setup](01-welcome-setup.md) (Pavol, 5 min)
1. [Welcome & Setup](02-tracing-introduction.md) (Matej, 10 min)
1. [Auto-instrumentation](03-auto-instrumentation.md) (Pavol & Anthony, 20 min)
1. [Manual-instrumentation](04-manual-instrumentation.md) (Bene & Matej, 10 min)
1. Wrap up & Questions
