# kBugger
[![Deploy and Release](https://github.com/JeffNeff/kBugger/actions/workflows/main.yml/badge.svg)](https://github.com/JeffNeff/kBugger/actions/workflows/main.yml)

`kBugger` is a knative service that allows the user to emit dynamic cloudevent types. 
These events will be emitted on a loop for the duration of the deployments lifecycle
with the enviorment variable `TIMEOUT` defining the sleep duration between requests. 

**Note** if `TIMEOUT` is set to 0 the service will emit a single event and then close

## Usage 
`kBugger` is basically just a cron job that allows you to specify cloudevent attributes. 
This service is especially useful when testing and debugging Knative Targets & Bridges.

To use `kBugger` in its simplest mannor, there is an example manifest located at `./kBugger.yaml`.
Included in this manifest is not only an example deployment of `kBugger` but also an
event display with its corresponding trigger implementation. 

Before deploying you will need to update at least one parameter in the manifest, `K_SINK`.
To do this, first decide on the namespace you will be deploying under (for this example `test` will be used),
and then retrieve the broker you will be using to deliver events.
```
$ kubectl -n test get brokers
NAME      URL                                                                     AGE   READY   REASON
default   http://broker-ingress.knative-eventing.svc.cluster.local/test/default   37d   True    
```

The URL here is the broker that we will be using to deliver our events. 

**note:** If there is not a Broker available you can create the default one by executing the following command from the root directory of this project:
```
kubectl -n test apply -f broker.yaml
```

After updating the `K_SINK` enviorment variable in the manifest, apply it with:
```
kubectl -n test apply -f kBugger.yaml
```

Verify that the service has been deployed:
```
kubectl -n test get pods
NAME                                                              READY   STATUS                       RESTARTS   AGE
alldisplay-00001-deployment-6467c885db-ftkhs                      2/2     Running                      0          7s
kbugger-00001-deployment-69978959bb-82ng9                         1/2     Running                      0          7s
```

Verify the service is working properly by checking the `alldisplay` logs
```
kubectl -n test logs alldisplay-00001-deployment-6467c885db-ftkhs user-container
Context Attributes,
  specversion: 1.0
  type: kbugger.subject
  source: kbugger
  subject: kbugger.subject
  id: io.kbugger.test
  time: 2021-07-09T19:25:07.214001959Z
  datacontenttype: application/json
Extensions,
  knativearrivaltime: 2021-07-09T19:25:07.214600301Z
  traceparent: 00-5fa955cca6dd2388a4a37c30bc74b89f-d00563bfaae89748-00
Data,
  {
    "test": "data"
  }
```



## Running locally

Export the following variables:
```
export EVENT_ID=
export EVENT_SOURCE=
export EVENT_SUBJECT=
export EVENT_TYPE=
export EVENT_DATA='{"test":"data"}'
export TIMEOUT=20
export K_SINK=localhost:8080
```

Run the project
```
cd cmd/kBugger
go run . 
```
