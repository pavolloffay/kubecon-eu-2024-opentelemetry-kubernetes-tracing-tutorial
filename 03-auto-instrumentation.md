# Auto-instrumentation

<!---  Recording https://vimeo.com/917479681/f57ec76c52   -->

## Application Description

The sample application is a simple _"dice game"_, where two players roll a
dice, and the player with the highest number wins.

There are 3 microservices within this application:

- Service `frontend` in Node.JS, that has an API endpoint `/` which takes two
  player names as query parameters (player1 and player2). The service calls 2
  down stream services (backend1, backend2), which each returning a random number
  between 1-6. The winner is computed and returned.
- Service `backend1` in python, that has an API endpoint `/rolldice` which takes
  a player name as query parameter. The service returns a random number between
  1 and 6.
- Service `backend2` in Java, that also has an API endpoint `/rolldice` which
  takes a player name as query parameter. The service returns a random number
  between 1 and 6.

Additionally there is a `loadgen` service, which utilizes `curl` to periodically
call the frontend service.

Let's assume player `alice` and `bob` use our service, here's a potential
sequence diagram:

```mermaid
sequenceDiagram
    loadgen->>frontend: /?player1=bob&player2=alice
    frontend->>backend1: /rolldice?player=bob
    frontend->>backend2: /rolldice?player=alice
    backend1-->>frontend: 3
    frontend-->>loadgen: bob rolls: 3
    backend2-->>frontend: 6
    frontend-->>loadgen: alice rolls: 6
    frontend-->>loadgen: alice wins
```

### Deploy the app into Kubernetes

Deploy the application into the kubernetes cluster. The app will be deployed into `tutorial-application` namespace.

```bash
kubectl apply -f https://raw.githubusercontent.com/pavolloffay/kubecon-eu-2024-opentelemetry-kubernetes-tracing-tutorial/main/app/k8s.yaml
kubectl get pods -n tutorial-application -w
...
NAME                                   READY   STATUS    RESTARTS   AGE
backend1-deployment-577cf945b4-tz5kv   1/1     Running   0          62s
backend2-deployment-59d4b47774-xbq84   1/1     Running   0          62s
frontend-deployment-678795956d-zwg4q   1/1     Running   0          62s
loadgen-deployment-5c7d6896f8-2fz6h    1/1     Running   0          62s
```

Now port-forward the frontend app:

```bash
kubectl port-forward service/frontend-service -n tutorial-application 4000:4000 
```

Open browser at [http://localhost:4000/](http://localhost:4000/).

## Manual or Automatic Instrumentation?

To make your application emit traces, metrics & logs you can either instrument
your application _manually_ or _automatically_:

- Manual instrumentation means that you modify your code yourself: you initialize and
  configure the SDK, you load instrumentation libraries, you create your own spans,
  metrics using the API.
  Developers can use this approach to tune the observability of their application to
  their needs, but it requires a lot of initial time investment, expertise how (RPC) frameworks and client work and maintenance over time.
- Automatic instrumentation means that you don't have to touch your code to get your
  application emit telemetry data.
  Automatic instrumentation is great to get you started with OpenTelemetry, and it is
  also valuable for Application Operators, who have no access or insights about the
  source code.

In this chapter we will cover using OpenTelemetry auto-instrumentation.

## Instrument the demo application

In this section we will deploy the app into Kubernetes and instrument it with OpenTelemetry auto-instrumentation
using the [Instrumentation CRD](https://github.com/open-telemetry/opentelemetry-operator?tab=readme-ov-file#opentelemetry-auto-instrumentation-injection) provided by the OpenTelemetry operator.
Then we will modify the app to create custom spans and collector additional attributes.

### Deploy OpenTelemetry collector

![OpenTelemetry Collector](images/otel-collector.png)

Deploy OpenTelemetry collector that will receive data from the instrumented workloads.

See the [OpenTelemetryCollector CR](./backend/03-collector.yaml).

```bash
kubectl apply -f https://raw.githubusercontent.com/pavolloffay/kubecon-eu-2024-opentelemetry-kubernetes-tracing-tutorial/main/backend/03-collector.yaml
kubectl get pods -n observability-backend -w
```

### Create instrumentation CR and see traces in the Jaeger console

Now let's instrument the app with the `Instrumentation` CR and see traces in the Jaeger console.

First the Instrumentation CR needs to be created in the `tutorial-application` namespace:

See the [Instrumentation CR](./app/instrumentation.yaml).

```bash
kubectl apply -f https://raw.githubusercontent.com/pavolloffay/kubecon-eu-2024-opentelemetry-kubernetes-tracing-tutorial/main/app/instrumentation.yaml
kubectl get pods -n tutorial-application -w
...                                                                                                                                                                                                                                                                                        
NAME                                   READY   STATUS    RESTARTS   AGE
backend1-deployment-577cf945b4-tz5kv   1/1     Running   0          8m59s
backend2-deployment-59d4b47774-xbq84   1/1     Running   0          8m59s
frontend-deployment-678795956d-zwg4q   1/1     Running   0          8m59s
loadgen-deployment-5c7d6896f8-2fz6h    1/1     Running   0          8m59s
```

The `Instrumentation` CR does not instrument the workloads. The instrumentation needs to be enabled by annotating a pod:

```bash
kubectl patch deployment frontend-deployment -n tutorial-application -p '{"spec": {"template":{"metadata":{"annotations":{"instrumentation.opentelemetry.io/inject-sdk":"true"}}}} }'
kubectl patch deployment backend1-deployment -n tutorial-application -p '{"spec": {"template":{"metadata":{"annotations":{"instrumentation.opentelemetry.io/inject-python":"true"}}}} }'
kubectl patch deployment backend2-deployment -n tutorial-application -p '{"spec": {"template":{"metadata":{"annotations":{"instrumentation.opentelemetry.io/inject-java":"true"}}}} }'
kubectl get pods -n tutorial-application -w
# Port forward again -> kubectl port-forward service/frontend-service -n tutorial-application 4000:4000 
...
NAME                                   READY   STATUS              RESTARTS   AGE
backend1-deployment-559946d88-c6zq7    0/1     Init:0/1            0          1s
backend2-deployment-5658ddfd6d-gz6ql   0/1     Init:0/1            0          1s
frontend-deployment-79b9c46d76-n74gr   0/1     ContainerCreating   0          1s
```

See the `backend2` pod spec:

```bash
kubectl describe pod backend2-deployment-5658ddfd6d-gz6ql -n tutorial-application
...
Init Containers:
  opentelemetry-auto-instrumentation-java:
    Image:         ghcr.io/open-telemetry/opentelemetry-operator/autoinstrumentation-java:1.32.1
    Command:
      cp
      /javaagent.jar
      /otel-auto-instrumentation-java/javaagent.jar
    Mounts:
      /otel-auto-instrumentation-java from opentelemetry-auto-instrumentation-java (rw)
      /var/run/secrets/kubernetes.io/serviceaccount from kube-api-access-48z6x (ro)
Containers:
  backend2:
    Image:          ghcr.io/pavolloffay/kubecon-eu-2024-opentelemetry-kubernetes-tracing-tutorial-backend2:latest
    Environment:
      OTEL_LOGS_EXPORTER:                  otlp
      JAVA_TOOL_OPTIONS:                    -javaagent:/otel-auto-instrumentation-java/javaagent.jar
      OTEL_SERVICE_NAME:                   backend2-deployment
      OTEL_EXPORTER_OTLP_ENDPOINT:         http://otel-collector.observability-backend.svc.cluster.local:4317
      OTEL_RESOURCE_ATTRIBUTES_POD_NAME:   backend2-deployment-5658ddfd6d-gz6ql (v1:metadata.name)
      OTEL_RESOURCE_ATTRIBUTES_NODE_NAME:   (v1:spec.nodeName)
      OTEL_PROPAGATORS:                    tracecontext,baggage,b3
      OTEL_TRACES_SAMPLER:                 parentbased_traceidratio
      OTEL_TRACES_SAMPLER_ARG:             1
      OTEL_RESOURCE_ATTRIBUTES:            k8s.container.name=backend2,k8s.deployment.name=backend2-deployment,k8s.namespace.name=tutorial-application,k8s.node.name=$(OTEL_RESOURCE_ATTRIBUTES_NODE_NAME),k8s.pod.name=$(OTEL_RESOURCE_ATTRIBUTES_POD_NAME),k8s.replicaset.name=backend2-deployment-5658ddfd6d,service.version=latest
    Mounts:
      /otel-auto-instrumentation-java from opentelemetry-auto-instrumentation-java (rw)
```

Now let's execute some requests on the app [http://localhost:4000/](http://localhost:4000/) and see traces in the Jaeger console [http://localhost:16686/](http://localhost:16686/).

![Trace search](./images/jaeger-trace-search.jpg)
![Trace detail](./images/jaeger-trace-detail.jpg)

In addition to traces in the Java auto-instrumentation also emits **logs** and **metrics**.
The logs in our case are printed into the collector stdout via `debug` exporter and metrics are sent via OTLP HTTP into Prometheus.
The OpenTelemetry spec defines that the following metrics should be collected: [HTTP metrics](https://opentelemetry.io/docs/specs/semconv/http/http-metrics/).

```bash
kubectl logs deployment.apps/otel-collector -n observability-backend
...
2024-02-28T10:08:21.807Z	info	LogsExporter	{"kind": "exporter", "data_type": "logs", "name": "debug", "resource logs": 1, "log records": 7}
2024-02-28T10:08:21.807Z	info	ResourceLog #0
Resource SchemaURL: https://opentelemetry.io/schemas/1.21.0
Resource attributes:
     -> container.id: Str(462d8e356c9b801d76edab5886730965f7f37b3d8b47d5eadfaea134141a35c1)
     -> host.arch: Str(amd64)
     -> host.name: Str(backend2-deployment-c7c8dc78c-wvhnk)
     -> k8s.container.name: Str(backend2)
     -> k8s.deployment.name: Str(backend2-deployment)
     -> k8s.namespace.name: Str(tutorial-application)
     -> k8s.node.name: Str(minikube)
     -> k8s.pod.name: Str(backend2-deployment-c7c8dc78c-wvhnk)
     -> k8s.replicaset.name: Str(backend2-deployment-c7c8dc78c)
     -> os.description: Str(Linux 6.5.12-100.fc37.x86_64)
     -> os.type: Str(linux)
     -> process.command_args: Slice(["/opt/java/openjdk/bin/java","-jar","./build/libs/dice-0.0.1-SNAPSHOT.jar"])
     -> process.executable.path: Str(/opt/java/openjdk/bin/java)
     -> process.pid: Int(7)
     -> process.runtime.description: Str(Eclipse Adoptium OpenJDK 64-Bit Server VM 21.0.2+13-LTS)
     -> process.runtime.name: Str(OpenJDK Runtime Environment)
     -> process.runtime.version: Str(21.0.2+13-LTS)
     -> service.name: Str(backend2-deployment)
     -> service.version: Str(withspan)
     -> telemetry.auto.version: Str(1.32.1)
     -> telemetry.sdk.language: Str(java)
     -> telemetry.sdk.name: Str(opentelemetry)
     -> telemetry.sdk.version: Str(1.34.1)
ScopeLogs #0
ScopeLogs SchemaURL: 
InstrumentationScope org.apache.catalina.core.ContainerBase.[Tomcat].[localhost].[/] 
LogRecord #0
ObservedTimestamp: 2024-02-28 10:08:21.178481174 +0000 UTC
Timestamp: 2024-02-28 10:08:21.178 +0000 UTC
SeverityText: INFO
SeverityNumber: Info(9)
Body: Str(Initializing Spring embedded WebApplicationContext)
Trace ID: 3bde5d3ee82303571bba6e1136781fe4 
Span ID: 45de5d3ee82303571bba6e1136781fe4
Flags: 0
ScopeLogs #1
ScopeLogs SchemaURL: 
InstrumentationScope io.opentelemetry.dice.DiceApplication 
LogRecord #0
ObservedTimestamp: 2024-02-28 10:08:21.638118261 +0000 UTC
Timestamp: 2024-02-28 10:08:21.638 +0000 UTC
SeverityText: INFO
SeverityNumber: Info(9)
Body: Str(Started DiceApplication in 3.459 seconds (process running for 6.305))
Trace ID: 3bde5d3ee82303571bba6e1136781fe4
Span ID: 46de5d3ee82303571bba6e1136781fe4
Flags: 0


kubectl logs -n tutorial-application deployment.apps/backend2-deployment
...
Defaulted container "backend2" out of: backend2, opentelemetry-auto-instrumentation-java (init)
Picked up JAVA_TOOL_OPTIONS:  -javaagent:/otel-auto-instrumentation-java/javaagent.jar
OpenJDK 64-Bit Server VM warning: Sharing is only supported for boot loader classes because bootstrap classpath has been appended
[otel.javaagent 2024-03-12 17:35:52:181 +0000] [main] INFO io.opentelemetry.javaagent.tooling.VersionLogger - opentelemetry-javaagent - version: 1.32.1

  .   ____          _            __ _ _
 /\\ / ___'_ __ _ _(_)_ __  __ _ \ \ \ \
( ( )\___ | '_ | '_| | '_ \/ _` | \ \ \ \
 \\/  ___)| |_)| | | | | || (_| |  ) ) ) )
  '  |____| .__|_| |_|_| |_\__, | / / / /
 =========|_|==============|___/=/_/_/_/
 :: Spring Boot ::                (v3.0.5)

2024-03-12T17:35:55.712Z  INFO 7 --- [           main] io.opentelemetry.dice.DiceApplication    : Starting DiceApplication v0.0.1-SNAPSHOT using Java 21.0.2 with PID 7 (/usr/src/app/build/libs/dice-0.0.1-SNAPSHOT.jar started by root in /usr/src/app)
2024-03-12T17:35:55.749Z  INFO 7 --- [           main] io.opentelemetry.dice.DiceApplication    : No active profile set, falling back to 1 default profile: "default"
2024-03-12T17:35:57.556Z  INFO 7 --- [           main] o.s.b.w.embedded.tomcat.TomcatWebServer  : Tomcat initialized with port(s): 5165 (http)
2024-03-12T17:35:57.588Z  INFO 7 --- [           main] o.apache.catalina.core.StandardService   : Starting service [Tomcat]
2024-03-12T17:35:57.589Z  INFO 7 --- [           main] o.apache.catalina.core.StandardEngine    : Starting Servlet engine: [Apache Tomcat/10.1.7]
2024-03-12T17:35:57.667Z  INFO 7 --- [           main] o.a.c.c.C.[Tomcat].[localhost].[/]       : Initializing Spring embedded WebApplicationContext
2024-03-12T17:35:57.669Z  INFO 7 --- [           main] w.s.c.ServletWebServerApplicationContext : Root WebApplicationContext: initialization completed in 1800 ms
2024-03-12T17:35:58.293Z  INFO 7 --- [           main] o.s.b.w.embedded.tomcat.TomcatWebServer  : Tomcat started on port(s): 5165 (http) with context path ''
2024-03-12T17:35:58.308Z  INFO 7 --- [           main] io.opentelemetry.dice.DiceApplication    : Started DiceApplication in 3.459 seconds (process running for 6.305)
2024-03-12T17:37:04.363Z  INFO 7 --- [nio-5165-exec-1] o.a.c.c.C.[Tomcat].[localhost].[/]       : Initializing Spring DispatcherServlet 'dispatcherServlet'
2024-03-12T17:37:04.364Z  INFO 7 --- [nio-5165-exec-1] o.s.web.servlet.DispatcherServlet        : Initializing Servlet 'dispatcherServlet'
2024-03-12T17:37:04.365Z  INFO 7 --- [nio-5165-exec-1] o.s.web.servlet.DispatcherServlet        : Completed initialization in 1 ms
2024-03-12T17:37:04.435Z  INFO 7 --- [nio-5165-exec-1] io.opentelemetry.dice.RollController     : Player 2 is rolling the dice: 2
2024-03-12T17:37:04.736Z  WARN 7 --- [nio-5165-exec-3] io.opentelemetry.dice.RollController     : Illegal number rolled, setting result to '1'
2024-03-12T17:37:04.737Z  INFO 7 --- [nio-5165-exec-3] io.opentelemetry.dice.RollController     : Player 2 is rolling the dice: 1
```

```bash
kubectl port-forward -n observability-backend service/prometheus 8080:80
```
Open Prometheus in the browser [localhost:8080](http://localhost:8080/graph?g0.expr=group%20(%7Bjob%3D%22backend2-deployment%22%7D)%20by%20(__name__)%0A&g0.tab=0&g0.stacked=0&g0.show_exemplars=0&g0.range_input=1h)

![Metrics from Java agent from backend2-deployment](./images/prometheus_javaagent_metrics_list.jpg)

### Customize Java auto-instrumentation with config (capture more data)

In this section we will configure the Java auto-instrumentation by modifying `Instrumentation` CR to:
* create custom spans - for the main method of the application
* capture server response HTTP headers

See the [Java agent docs](https://opentelemetry.io/docs/languages/java/automatic/configuration/) with all the configuration options.

See the [Instrumentation CR](./app/instrumentation-java-custom-config.yaml).

```bash
kubectl apply -f https://raw.githubusercontent.com/pavolloffay/kubecon-eu-2024-opentelemetry-kubernetes-tracing-tutorial/main/app/instrumentation-java-custom-config.yaml
kubectl rollout restart deployment.apps/backend2-deployment -n tutorial-application
kubectl get pods -w -n tutorial-application
```

![Span from backend2-deployment](./images/jaeger-capture-custom-headers.jpg)

### Customize Java auto-instrumentation with code (capture more data)

> [!NOTE]  
> This is an optional more advanced section.

In this section we will modify [Java backend2](./app/backend2) service to:
* create a new span to observe execution of a business method
* attach attributes to span

The OpenTelemetry Java auto-instrumentation supports `@WithSpan`, `@SpanAttribute` and `@AddingSpanAttributes` see the [documentation](https://opentelemetry.io/docs/languages/java/automatic/annotations/) and [javadoc](https://javadoc.io/doc/io.opentelemetry.instrumentation/opentelemetry-instrumentation-annotations/latest/io/opentelemetry/instrumentation/annotations/package-summary.html).

Open the [RollController.java](./app/backend2/src/main/java/io/opentelemetry/dice/RollController.java) and use the annotations:

```java
# app/backend2/build.gradle
#   implementation 'io.opentelemetry.instrumentation:opentelemetry-instrumentation-annotations:2.1.0'
#   implementation 'io.opentelemetry:opentelemetry-api:1.35.0'
 
    import io.opentelemetry.api.trace.Span;
    import io.opentelemetry.instrumentation.annotations.WithSpan;
    import io.opentelemetry.instrumentation.annotations.SpanAttribute;
    import io.opentelemetry.instrumentation.annotations.AddingSpanAttributes;

    @AddingSpanAttributes
	@GetMapping("/rolldice")
	public String index(@SpanAttribute("player") @RequestParam("player") Optional<String> player) {
    
    @WithSpan
    public int getRandomNumber(@SpanAttribute("min") int min, @SpanAttribute("max") int max) {
        int result = (int) ((Math.random() * (max - min)) + min);
        Span span = Span.current();
        span.setAttribute("result", result);
        return result;
    }
```

Compile it and deploy:
```bash
cd app/backend2

# Use minikube's docker registry
# eval $(minikube -p minikube docker-env)
docker build -t ghcr.io/pavolloffay/kubecon-eu-2024-opentelemetry-kubernetes-tracing-tutorial-backend2:withspan . 
# docker push ghcr.io/pavolloffay/kubecon-eu-2024-opentelemetry-kubernetes-tracing-tutorial-backend2:withspan

kubectl set image deployment.apps/backend2-deployment backend2=ghcr.io/pavolloffay/kubecon-eu-2024-opentelemetry-kubernetes-tracing-tutorial-backend2:withspan -n tutorial-application
kubectl get pods -w -n tutorial-application
```

![Span from backend2-deployment](./images/jaeger-with-span.jpg)

---
[Next steps](./04-manual-instrumentation.md)
