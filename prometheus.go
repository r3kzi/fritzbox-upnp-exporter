package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
)

type FritzBoxCollector struct {
	Config *Config
}

type Metric struct {
	Name string
	Desc *prometheus.Desc
}

var metrics = []Metric{
	{
		Name: "PhysicalLinkStatus",
		Desc: prometheus.NewDesc(
			"fritzbox_wan_layer1_link_status",
			"Status of physical link (Up = 1)",
			nil,
			nil,
		),
	},
	{
		Name: "Layer1DownstreamMaxBitRate",
		Desc: prometheus.NewDesc(
			"fritzbox_wan_layer1_downstream_max_bitrate",
			"Layer1 downstream max bitrate",
			nil,
			nil,
		),
	},
	{
		Name: "Layer1UpstreamMaxBitRate",
		Desc: prometheus.NewDesc(
			"fritzbox_wan_layer1_upstream_max_bitrate",
			"Layer1 upstream max bitrate",
			nil,
			nil,
		),
	},
	{
		Name: "TotalBytesSent",
		Desc: prometheus.NewDesc(
			"fritzbox_wan_bytes_sent",
			"bytes sent on gateway WAN interface",
			nil,
			nil,
		),
	},
	{
		Name: "TotalBytesReceived",
		Desc: prometheus.NewDesc(
			"fritzbox_wan_bytes_received",
			"bytes received on gateway WAN interface",
			nil,
			nil,
		),
	},
	{
		Name: "TotalPacketsSent",
		Desc: prometheus.NewDesc(
			"fritzbox_wan_packets_sent",
			"packets sent on gateway WAN interface",
			nil,
			nil,
		),
	},
	{
		Name: "TotalPacketsReceived",
		Desc: prometheus.NewDesc(
			"fritzbox_wan_packets_received",
			"packets received on gateway WAN interface",
			nil,
			nil,
		),
	},
}

func newFritzBoxCollector(config *Config) *FritzBoxCollector {
	return &FritzBoxCollector{
		Config: config,
	}
}

func (collector *FritzBoxCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range metrics {
		ch <- metric.Desc
	}
}

func (collector *FritzBoxCollector) Collect(ch chan<- prometheus.Metric) {
	uPnPClient := NewUPnPClient(collector.Config.URL, collector.Config.Username, collector.Config.Password)
	values := uPnPClient.Execute()

	extract := func(metric string) float64 {
		for k, v := range values {
			if k == metric {
				if s, err := strconv.ParseFloat(v, 64); err == nil {
					return s
				} else {
					if v == "Up" {
						return 1.0
					}
				}
			}
		}
		return 0.0
	}

	for _, metric := range metrics {
		ch <- prometheus.MustNewConstMetric(metric.Desc, prometheus.GaugeValue, extract(metric.Name))
	}
}
