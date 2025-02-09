# ExternalDNS - Technitium DNS Webhook

ExternalDNS is a Kubernetes add-on for automatically managing
Domain Name System (DNS) records for Kubernetes services by using different DNS providers.
The [Technitium DNS](https://technitium.com/dns/) webhook allows to manage your
TechnitiumDNS domains inside your kubernetes cluster with [ExternalDNS](https://github.com/kubernetes-sigs/external-dns).

As with many other ExternalDNS webhook providers,
this project is heavily inspired and helped by the [Ionos External DNS](https://github.com/ionos-cloud/external-dns-ionos-webhook)
webhook implementation.

## Kubernetes Deployment

```shell
helm repo add external-dns https://kubernetes-sigs.github.io/external-dns/

kubectl create secret generic technitium-credentials --from-literal=user='<USER>' --from-literal=pass='<PASS>'

# create the helm values file
cat <<EOF > external-dns-technitium-values.yaml
image:
  tag: v0.15.0

# -- ExternalDNS Log level.
logLevel: debug # reduce in production

# -- if true, _ExternalDNS_ will run in a namespaced scope (Role and Rolebinding will be namespaced too).
namespaced: false

# -- _Kubernetes_ resources to monitor for DNS entries.
sources:
  - ingress
  - service
  - crd

extraArgs:
  ## must override the default value with port 8888 with port 8080 because this is hard-coded in the helm chart
  - --webhook-provider-url=http://localhost:8080

provider:
  name: webhook
  webhook:
    image:
      repository: ghcr.io/chrisatcho/external-dns-technitium-webhook
      tag: v0.0.1
    env:
    - name: LOG_LEVEL
      value: debug # reduce in production
    - name: TECHNITIUM_API_URL
      value: "http://technitiumdns:5380"
    - name: TECHNITIUM_USER
      valueFrom:
        secretKeyRef:
          name: technitium-credentials
          key: user
    - name: TECHNITIUM_PASS
      valueFrom:
        secretKeyRef:
          name: technitium-credentials
          key: pass
    - name: SERVER_HOST
      value: "0.0.0.0"
    - name: SERVER_PORT
      value: "8080"
    - name: TECHNITIUM_DEBUG
      value: "false" # put this to true if you want see details of the http requests
    livenessProbe:
      httpGet:
        path: /health
    readinessProbe:
      httpGet:
        path: /health
EOF

# install external-dns with helm
helm upgrade external-dns-technitium external-dns/external-dns --version 1.15.0 -f external-dns-technitium-values.yaml --install
```

## Configuration

### Technitium Configuration

| Environment Variable | Description                  | Default |
| -------------------- | ---------------------------- | ------- |
| `TECHNITIUM_USER`    | Username                     | None    |
| `TECHNITIUM_PASS`    | Password                     | None    |
| `TECHNITIUM_API_URL` | Full url of the API endpoint | None    |
| `TECHNITIUM_DEBUG`   | Enable / Disable API logging | `False` |

### Server Configuration

| Environment Variable             | Description                                                      | Default Value |
| -------------------------------- | ---------------------------------------------------------------- | ------------- |
| `SERVER_HOST`                    | The host address where the server listens.                       | `localhost`   |
| `SERVER_PORT`                    | The port where the server listens.                               | `8888`        |
| `SERVER_READ_TIMEOUT`            | Duration the server waits before timing out on read operations.  | N/A           |
| `SERVER_WRITE_TIMEOUT`           | Duration the server waits before timing out on write operations. | N/A           |
| `DOMAIN_FILTER`                  | List of domains to include in the filter.                        | Empty         |
| `EXCLUDE_DOMAIN_FILTER`          | List of domains to exclude from filtering.                       | Empty         |
| `REGEXP_DOMAIN_FILTER`           | Regular expression for filtering domains.                        | Empty         |
| `REGEXP_DOMAIN_FILTER_EXCLUSION` | Regular expression for excluding domains from the filter.        | Empty         |
