# FritzBox UPnP Prometheus Exporter

- build on top of [TR064](https://avm.de/service/schnittstellen/)
- uses HTTP Digest Authentication with TLS

## Currently exposed metrics

- `WANAccessType`
- `Layer1UpstreamMaxBitRate`
- `Layer1DownstreamMaxBitRate`
- `PhysicalLinkStatus`
- `TotalBytesSent`
- `TotalBytesReceived`
- `TotalPacketsSent`
- `TotalPacketsReceived`

```
$ curl localhost:8080/metrics

# HELP fritzbox_wan_bytes_received bytes received on gateway WAN interface
# TYPE fritzbox_wan_bytes_received gauge
fritzbox_wan_bytes_received 1.960837653e+09
# HELP fritzbox_wan_bytes_sent bytes sent on gateway WAN interface
# TYPE fritzbox_wan_bytes_sent gauge
fritzbox_wan_bytes_sent 2.34091343e+08
# HELP fritzbox_wan_layer1_downstream_max_bitrate Layer1 downstream max bitrate
# TYPE fritzbox_wan_layer1_downstream_max_bitrate gauge
fritzbox_wan_layer1_downstream_max_bitrate 100000
# HELP fritzbox_wan_layer1_link_status Status of physical link (Up = 1)
# TYPE fritzbox_wan_layer1_link_status gauge
fritzbox_wan_layer1_link_status 1
# HELP fritzbox_wan_layer1_upstream_max_bitrate Layer1 upstream max bitrate
# TYPE fritzbox_wan_layer1_upstream_max_bitrate gauge
fritzbox_wan_layer1_upstream_max_bitrate 50000
# HELP fritzbox_wan_packets_received packets received on gateway WAN interface
# TYPE fritzbox_wan_packets_received gauge
fritzbox_wan_packets_received 505030
# HELP fritzbox_wan_packets_sent packets sent on gateway WAN interface
# TYPE fritzbox_wan_packets_sent gauge
fritzbox_wan_packets_sent 803660
```