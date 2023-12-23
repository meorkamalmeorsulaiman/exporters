## This exporters use json to pull interface tx/rx bytes

### Docker images

- Below are the list of docker images used for this lab:

```bash
kamal.sulaiman@kamal exporters % docker images
REPOSITORY                          TAG       IMAGE ID       CREATED        SIZE
grafana/grafana                     latest    05e0ef008060   4 days ago     392MB
prom/prometheus                     latest    929804fcacca   2 weeks ago    242MB
prometheuscommunity/json-exporter   latest    399788366d89   6 months ago   22.1MB
```

- Spining it up, this were run in individual container, you may use any kind to make it as a service

```bash
docker run -d --name=grafana -p 3000:3000 grafana/grafana
docker run -d --name prom -v ./prometheus.yml:/etc/prometheus/prometheus.yml -p 9090:9090 prom/prometheus
docker run -d --name json -v ./config_mist.yml:/config.yml -p 7979:7979 prometheuscommunity/json-exporter
```

### JSON Exporter configs

- Below is the json-exporter configurations, file included `config_mist.yml`

```yaml
---
modules:
  if_stats:
    metrics:
    - name: if_transmit
      help: Output bytes
      type: object
      path: '{ .results[*] }'
      labels:
        ifname: '{ .port_id }'
      valuetype: counter
      values:
        bytes: '{ .tx_bytes }'

    - name: if_received
      help: Input bytes
      type: object
      path: '{ .results[*] }'
      labels:
        ifname: '{ .port_id }'
      valuetype: counter
      values:
        bytes: '{ .rx_bytes }'

    http_client_config:
      tls_config:
        insecure_skip_verify: true
      authorization:
        type: Token
        credentials: ${TOKEN}
```

### Prometheus configs

- API docs can access thru `api.mist.com` and files included `prometheus.yml`

```yaml
scrape_configs:

  ## gather metrics of prometheus itself
- job_name: prometheus
  static_configs:
    - targets:
      - localhost:9090 # equivalent to "localhost:9090"

  ## gather the metrics of json_exporter application itself
# - job_name: json_exporter
#   static_configs:
#     - targets:
#       ## Location of the json exporter's real <hostname>:<port>
#       - localhost:7979 # equivalent to "localhost:7979"

- job_name: json
  metrics_path: /probe
  params:
    module: [if_stats]
  static_configs:
    - targets:
      - https://api.mist.com/api/v1/sites/:SITE/stats/switch_ports/search?mac=:MAC #One Target per switch
      - https://api.mist.com/api/v1/sites/:SITE/stats/switch_ports/search?mac=:MAC
      - https://api.mist.com/api/v1/sites/:SITE/stats/switch_ports/search?mac=:MAC
  relabel_configs:
    - source_labels: [__address__]
      target_label: __param_target
    - source_labels: [__param_target]
      target_label: instance
    - target_label: __address__
      ## Location of the json exporter's real <hostname>:<port>
      replacement: localhost:7979 # equivalent to "localhost:7979"

  ## gather the metrics from third party json sources, via the json exporter
```

### Sample metrics scrape from API

```
# HELP if_received_bytes Input bytes
# TYPE if_received_bytes counter
if_received_bytes{ifname="ae0"} 2.5069269271982e+13
if_received_bytes{ifname="ae101"} 2.021562381275e+12
if_received_bytes{ifname="ae102"} 3.772352658277e+12
if_received_bytes{ifname="ae103"} 5.52152065091e+11
if_received_bytes{ifname="ae104"} 4.812101479556e+12
if_received_bytes{ifname="ae105"} 4.651374148677e+12
if_received_bytes{ifname="ae106"} 2.499506439405e+12
if_received_bytes{ifname="ae107"} 2.8430754977821e+13
if_received_bytes{ifname="ae108"} 3.8996374881686e+13
if_received_bytes{ifname="ae109"} 2.8279452340995e+13
if_received_bytes{ifname="ae110"} 2.0388835396758e+13
if_received_bytes{ifname="em0"} 3.8626855517e+10
if_received_bytes{ifname="em2"} 3.0796385105e+10
if_received_bytes{ifname="em3"} 3.08401542874e+11
if_received_bytes{ifname="ge-0/0/10"} 1.573004351642e+12
if_received_bytes{ifname="ge-0/0/11"} 1.822038116889e+12
if_received_bytes{ifname="vme"} 3.8371761324e+10
if_received_bytes{ifname="xe-0/0/0"} 2.021562380016e+12
if_received_bytes{ifname="xe-0/0/1"} 3.772352465408e+12
if_received_bytes{ifname="xe-0/0/12"} 2.8430752936289e+13
if_received_bytes{ifname="xe-0/0/13"} 3.8996370930901e+13
if_received_bytes{ifname="xe-0/0/14"} 2.8279450242439e+13
if_received_bytes{ifname="xe-0/0/15"} 2.038883317937e+13
if_received_bytes{ifname="xe-0/0/2"} 5.52152064014e+11
if_received_bytes{ifname="xe-0/0/3"} 4.812101478198e+12
if_received_bytes{ifname="xe-0/0/4"} 4.651374146831e+12
if_received_bytes{ifname="xe-0/0/44"} 6.747533115406e+12
if_received_bytes{ifname="xe-0/0/45"} 6.506242877092e+12
if_received_bytes{ifname="xe-0/0/46"} 1.2522399274498e+13
if_received_bytes{ifname="xe-0/0/47"} 1.2546869997484e+13
if_received_bytes{ifname="xe-0/0/5"} 2.499506438851e+12
# HELP if_transmit_bytes Output bytes
# TYPE if_transmit_bytes counter
if_transmit_bytes{ifname="ae0"} 4.4776052882281e+13
if_transmit_bytes{ifname="ae101"} 2.058975329948e+12
if_transmit_bytes{ifname="ae102"} 8.8243387546471e+13
if_transmit_bytes{ifname="ae103"} 5.05294818906e+11
if_transmit_bytes{ifname="ae104"} 4.121720845877e+12
if_transmit_bytes{ifname="ae105"} 3.106717243106e+12
if_transmit_bytes{ifname="ae106"} 1.14421038332e+12
if_transmit_bytes{ifname="ae107"} 4.448477465743e+12
if_transmit_bytes{ifname="ae108"} 5.598057151622e+12
if_transmit_bytes{ifname="ae109"} 7.017555192336e+12
if_transmit_bytes{ifname="ae110"} 2.077570725222e+12
if_transmit_bytes{ifname="bme0"} 4.4165724761e+10
if_transmit_bytes{ifname="em0"} 1.91544615843e+11
if_transmit_bytes{ifname="em2"} 2.680189499e+09
if_transmit_bytes{ifname="em3"} 4.5587783115e+10
if_transmit_bytes{ifname="ge-0/0/10"} 7.82545064953e+11
if_transmit_bytes{ifname="ge-0/0/11"} 2.159996156e+09
if_transmit_bytes{ifname="vme"} 1.91381235041e+11
if_transmit_bytes{ifname="xe-0/0/0"} 2.058975286431e+12
if_transmit_bytes{ifname="xe-0/0/1"} 8.8243371057504e+13
if_transmit_bytes{ifname="xe-0/0/12"} 4.448477465743e+12
if_transmit_bytes{ifname="xe-0/0/13"} 5.598057151622e+12
if_transmit_bytes{ifname="xe-0/0/14"} 7.017555192336e+12
if_transmit_bytes{ifname="xe-0/0/15"} 2.077570725222e+12
if_transmit_bytes{ifname="xe-0/0/2"} 5.05294776194e+11
if_transmit_bytes{ifname="xe-0/0/3"} 4.121720802727e+12
if_transmit_bytes{ifname="xe-0/0/4"} 3.106717201612e+12
if_transmit_bytes{ifname="xe-0/0/44"} 6.781414316367e+12
if_transmit_bytes{ifname="xe-0/0/45"} 7.059591172996e+12
if_transmit_bytes{ifname="xe-0/0/46"} 2.2520337330008e+13
if_transmit_bytes{ifname="xe-0/0/47"} 2.2255715552273e+13
if_transmit_bytes{ifname="xe-0/0/5"} 1.144210340608e+12
```