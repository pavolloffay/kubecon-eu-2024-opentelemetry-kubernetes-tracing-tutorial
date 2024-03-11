# OpenTelemetry Transformation Language and Spans (OTTL)

## Overview

The OpenTelemetry Transformation Language (OTTL) is a powerful language that allows you to transform telmetry data flowing through the Collector. The transformation of data is executed on the basis of OTTL statements, which define how the telemetry should be transformed. It is a stand-alone part of the Collector codebase that is (re)used in a number of components, such as `filterprocessor`, `transformprocessor` or `routingprocessor`. 

Statements follow the OTTL grammar and are defined in the configuration of the particular component. Statements always relates to a particular **context** and invoke specific **functions** on the context. Besides that, as with any programming language, you can use operators for comparing values, convereters or literals. Combined all together, an example statement could look like this:

```yaml
    set(attributes["client_error"], true) where attributes["http.status"] == 400 or attributes["http.status"] == 404
```

In this statement, we're using the `set` function to set the `client_error` attribute to `true`. We're conditioning this by using the `where` qualified, to apply only if HTTP status code of the requests is `400` OR `404`.

### Contexts

Context determines which part of telemetry data should the statement be applied to. This can be universal for all signals, such as `Resource` and `Instrumentation Scope`, or they differ depending on the type of the signal such `Span`, `Datapoint` or `Log`. In the statement, a particular part of the context can be accessed via [paths](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/ottl/LANGUAGE.md#paths), which support the familiar `.` notation and accessing particular keys with `[]` (see the example above - `attributes["client_error"]` - accessing a particular attribute). Full list of all contexts can be found [here](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/pkg/ottl#getting-started).

### Functions

OTTL provides a list of predefind functions that come in two flavors - **editors** and **converters**. Editors work directly on and transform telemetry itself. Editors function include functions such as a `set`, `delete_key`, `replace_match` or `limit`. Conversely, converters are used to transform input within a statement and they do **not** modify the telemetry themselves. These can be used e.g. to get input length (`Len`), manipulate strings (`Concat`) or assert types (`IsInt`, `IsMap`, `IsString`...). Full list of both types of functions can be found [here](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/pkg/ottl/ottlfuncs#ottl-functions). 

### Other language features (grammar)

As mentioned, OTTL also supports other language features such as literals, operators, and comments. Full list of these can be found [here](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/pkg/ottl/LANGUAGE.md#paths), but most of these are common to many programming / scripting langauges and are fairly intuitive.

### OTTL in action

Let's take a look how OTTL can look in action with tracing. OTTL can be leveraged to transform spans in a useful ways, whether you want to enrich your spans with extra data, remove sensitive information or limit the amount of metadata included with your spans. To do this, we will use the `transformprocessor`. 

Our application supports recording the name of the players who are rolling the dice, by passing the name of the platers as parameters in the URL, e.g. `?player1=John&player2=Jane`. Due to privacy concerns, we might not want to include these names as attributes on our spans and we would rather annonymize them. So to do this, we will always pick only the first letter of the players name and include it as the attribute.

First, take a look at Jaeger and see that our spans have the `app.player1` and `app.player2` attributes. Choose the `frontend-deployment` service and observe that the root span has attributes `app.player1` and `app.player2` with the full names of the players.

Second, inspect the configuration for our `transformprocessor` below:

```yaml
processors:
    transform:
        error_mode: ignore
        trace_statements:
            - context: span
            statements:
                - set(attributes["app.player1"], Substring(attributes["app.player1"], 0, 1)) where attributes["app.player1"] != ""
                - set(attributes["app.player2"], Substring(attributes["app.player2"], 0, 1)) where attributes["app.player2"] != ""
```

We're using the `span` context to change the desired attributes of our spans. We're going to look for the `app.player1` and `app.player2` attributes and set them to the first letter of the name. We're using the `Substring` editor to do this, which takes the string, the start position and the length of the substring. We're also using the `where` qualifier to only apply this transformation if the attribute is not empty.

But that is not everything. Do you see the `http.url` and `http.url` attributes? These attributes still include names of our players as URL parameters! We need to get rid of them here as well. To achieve this, we add one more statement to replace the player name with `{playerName}` placeholder.

```yaml
processors:
    transform:
        error_mode: ignore
        trace_statements:
            - context: span
            statements:
                - set(attributes["app.player1"], Substring(attributes["app.player1"], 0, 1)) where attributes["app.player1"] != ""
                - set(attributes["app.player2"], Substring(attributes["app.player2"], 0, 1)) where attributes["app.player2"] != ""
                - replace_all_patterns(attributes, "value", "player1=[a-zA-Z]*", "player1={playerName}")
                - replace_all_patterns(attributes, "value", "player2=[a-zA-Z]*", "player2={playerName}")

```

Apply the changes to our collector, with transform processor now enabled in our tracing pipeline:

```bash
kubectl replace -f https://raw.githubusercontent.com/pavolloffay/kubecon-eu-2024-opentelemetry-kubernetes-tracing-tutorial/main/backend/07-collector.yaml
kubectl get pods -n observability-backend -w
```

After the collector with new configuration rolls out, run a couple requests with player names set:
- http://localhost:4000/?player1=John_Doe&player2=Jane_Doe
- http://localhost:4000/?player1=Neo&player2=Trinity
- http://localhost:4000/?player1=Barbie&player2=Ken

Now open your [Jaeger UI](http://localhost:16686/) and observe the spans. You should see that the `app.player1` and `app.player2` attributes are now anonymized and the player names are now replaced with `{playerName}` in attributes that contain the URL. You have succesfully transformed your spans with OTTL!