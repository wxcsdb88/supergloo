# Ingress

## Open questions

- Security / certs

## Istio Configuration

See also: https://istio.io/docs/tasks/traffic-management/ingress/

Istio Gateway:

```yaml
apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  name: httpbin-gateway
spec:
  selector:
    istio: ingressgateway # use Istio default gateway implementation
  servers:
  - port:
      number: 80
      name: http
      protocol: HTTP
    hosts:
    - "httpbin.example.com"
```

Istio VirtualService:

```yaml
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: httpbin
spec:
  hosts:
  - "httpbin.example.com"
  gateways:
  - httpbin-gateway
  http:
  - match:
    - uri:
        prefix: /status
    - uri:
        prefix: /delay
    route:
    - destination:
        port:
          number: 8000
        host: httpbin
```

## Linkerd Configuration

See also: https://blog.buoyant.io/2017/04/06/a-service-mesh-for-kubernetes-part-viii-linkerd-as-an-ingress-controller/

Linkerd ingress config map, DaemonSet, Service

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: l5d-config
data:
  config.yaml: |-
    admin:
      ip: 0.0.0.0
      port: 9990

    namers:
    - kind: io.l5d.k8s

    routers:
    - protocol: http
      identifier:
        kind: io.l5d.ingress
      servers:
        - port: 80
          ip: 0.0.0.0
          clearContext: true
      dtab: /svc => /#/io.l5d.k8s

    usage:
      orgId: linkerd-examples-ingress
---
apiVersion: extensions/v1beta1
kind: DaemonSet
metadata:
  labels:
    app: l5d
  name: l5d
spec:
  template:
    metadata:
      labels:
        app: l5d
    spec:
      volumes:
      - name: l5d-config
        configMap:
          name: "l5d-config"
      containers:
      - name: l5d
        image: buoyantio/linkerd:1.4.6
        env:
        - name: POD_IP
          valueFrom:
            fieldRef:
              fieldPath: status.podIP
        args:
        - /io.buoyant/linkerd/config/config.yaml
        ports:
        - name: http
          containerPort: 80
        - name: admin
          containerPort: 9990
        volumeMounts:
        - name: "l5d-config"
          mountPath: "/io.buoyant/linkerd/config"
          readOnly: true

      - name: kubectl
        image: buoyantio/kubectl:v1.8.5
        args: ["proxy", "-p", "8001"]
---
apiVersion: v1
kind: Service
metadata:
  name: l5d
spec:
  selector:
    app: l5d
  type: LoadBalancer
  ports:
  - name: http
    port: 80
  - name: admin
    port: 9990

```

```yaml
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: hello-world
  annotations:
    kubernetes.io/ingress.class: "linkerd"
spec:
  backend:
    serviceName: world-v1
    servicePort: http
  rules:
  - host: world.v2
    http:
      paths:
      - backend:
          serviceName: world-v2
          servicePort: http
```
