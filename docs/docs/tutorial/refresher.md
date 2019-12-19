---
name: Refresher on deploying
menu: Tutorial
route: /tutorial/refresher
---

# Refresher on deploying

## Deploying to Kubernetes using kubectl

To understand how Tanka works, it is important to know what steps are required
for the task of deploying Grafana and Prometheus to Kubernetes:

1. Prometheus
   - A `Deployment` must be created, to run the `prom/prometheus` image
   - Also a `Service` is needed for Grafana to be able to connect port `9090` of
     Prometheus.
2. Grafana
   - Another `Deployment` is required for the Grafana server.
   - To connect to the webinterface, we also need a `Service` of type
     `NodePort`.

Before taking a look how Tanka can help doing so, let's recall how to do it with
plain `kubectl`.

## Writing the yaml

`kubectl` expects the resources it should create in `.yaml` format. For Grafana
...

##### grafana.yaml:

```yaml
# Grafana server Deployment
apiVersion: apps/v1
kind: Deployment
metadata:
  name: grafana
spec:
  selector:
    matchLabels:
      name: grafana
  template:
    metadata:
      labels:
        name: grafana
    spec:
      containers:
        - image: grafana/grafana
          name: grafana
          ports:
            - containerPort: 3000
              name: ui
---
# Grafana UI Service NodePort
apiVersion: v1
kind: Service
metadata:
  labels:
    name: grafana
  name: grafana
spec:
  ports:
    - name: grafana-ui
      port: 3000
      targetPort: 3000
  selector:
    name: grafana
  type: NodePort
```

... and for Prometheus:

##### prometheus.yaml

```yaml
# Prometheus server Deployment
apiVersion: apps/v1
kind: Deployment
metadata:
  name: prometheus
spec:
  selector:
    matchLabels:
      name: prometheus
  template:
    metadata:
      labels:
        name: prometheus
    spec:
      containers:
        - image: prom/prometheus
          name: prometheus
          ports:
            - containerPort: 9090
              name: api
---
# Prometheus API Service
apiVersion: v1
kind: Service
metadata:
  labels:
    name: prometheus
  name: prometheus
spec:
  ports:
    - name: prometheus-api
      port: 9090
      targetPort: 9090
  selector:
    name: prometheus
```

That's pretty verbose, right?

Even worse, there are labels and matchers (e.g. `prometheus`) that need to be
exactly the same scattered across the file. It's a nightmare to debug and
furthermore harms readability a lot.

## Deploying to the cluster

To actually apply those resources, copy them into `.yaml` files and use:

```bash
$ kubectl apply -f prometheus.yaml grafana.yaml
deployment.apps/grafana created
deployment.apps/prometheus created
service/grafana created
service/prometheus created
```

## Checking it worked

So far so good, but can we tell it actually did what we wanted? Let's test that
Grafana can connect to Prometheus!

```bash
# Temporarily forward Grafana to localhost
kubectl port-forward deployments/grafana 3000:8080
```

Now go to http://localhost:8080 in your browser and login using `admin:admin`.
Then navigate to `Configuration > Data Sources > Add data source`, choose
`Prometheus` as type and enter `http://prometheus:9090` as URL. Hit
`Save & Test` which should yield a big green bar telling you everything is good.

Cool! This worked out well for this small example, but the `.yaml` files are
hard to read and maintain. Especially when you need to deploy this exact same
thing in `dev` and `prod` your choices are very limited.

Let's explore how Tanka can help us here in the next section!

## Cleaning up

Let's remove everything we created to start fresh with Jsonnet in the next section:

```bash
$ kubectl delete -f prometheus.yaml grafana.yaml
```
