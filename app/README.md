# Sample Application

This is the source code of the sample application used in the tutorial step [Deploy the application](../03-app-instrumentation.md).

## Usage

### Docker Compose

The quickest way to try out the sample application is by using docker compose.
In the [app/](.) folder run:

```console
docker compose up
```

This will build the images for all services + the load generator and then run
them for you. You can test the frontend service yourself by calling it with
your browser or curl at <http://localhost:4000?player1=alice&player2=bob>

If you'd like to try out all services being instrumented with OpenTelemetry run

```console
docker compose --env-file=./otel-env up
```

### Kubernetes

To run the the sample application on your kubernetes cluster, run

```console
kubectl apply -f ./k8s.yaml
```

Note that this will pull images of the applications from ghcr.io.

If you'd like to access the frontend service, open a new terminal and run

```console
kubectl port-forward svc/frontend-service 4000:4000
```

### Development

If you'd like to change the code of any of the applications,
you need to install Node.JS, Java and python3 with flask first.

Then you can run them all standalone:

- frontend (in folder [./frontend](./frontend)):
  
  ```console
  node index.js
  ```

- backend1 (in folder [./backend1](./backend1)):

  ```console
  flask run
  ```

- backend2 (in folder [./backend2](./backend)):

  ```console
   ./gradlew bootRun
  ```
