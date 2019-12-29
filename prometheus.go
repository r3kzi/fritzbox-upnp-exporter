package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"strconv"
)

type FritzBoxCollector struct {
	PhysicalLinkStatus         *prometheus.Desc
	Layer1DownstreamMaxBitRate *prometheus.Desc
	Layer1UpstreamMaxBitRate   *prometheus.Desc
	TotalBytesSent             *prometheus.Desc
	TotalBytesReceived         *prometheus.Desc
	TotalPacketsSent           *prometheus.Desc
	TotalPacketsReceived       *prometheus.Desc
}

func newFritzBoxCollector() *FritzBoxCollector {
	return &FritzBoxCollector{
		PhysicalLinkStatus: prometheus.NewDesc(
			"fritzbox_wan_layer1_link_status",
			"Status of physical link (Up = 1)",
			nil,
			nil,
		),
		Layer1DownstreamMaxBitRate: prometheus.NewDesc(
			"fritzbox_wan_layer1_downstream_max_bitrate",
			"Layer1 downstream max bitrate",
			nil,
			nil,
		),
		Layer1UpstreamMaxBitRate: prometheus.NewDesc(
			"fritzbox_wan_layer1_upstream_max_bitrate",
			"Layer1 upstream max bitrate",
			nil,
			nil,
		),
		TotalBytesSent: prometheus.NewDesc(
			"fritzbox_wan_bytes_sent",
			"bytes sent on gateway WAN interface",
			nil,
			nil,
		),
		TotalBytesReceived: prometheus.NewDesc(
			"fritzbox_wan_bytes_received",
			"bytes received on gateway WAN interface",
			nil,
			nil,
		),
		TotalPacketsSent: prometheus.NewDesc(
			"fritzbox_wan_packets_sent",
			"packets sent on gateway WAN interface",
			nil,
			nil,
		),
		TotalPacketsReceived: prometheus.NewDesc(
			"fritzbox_wan_packets_received",
			"packets received on gateway WAN interface",
			nil,
			nil,
		),
	}
}

func (collector *FritzBoxCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.PhysicalLinkStatus
	ch <- collector.Layer1DownstreamMaxBitRate
	ch <- collector.Layer1UpstreamMaxBitRate
	ch <- collector.TotalBytesSent
	ch <- collector.TotalBytesReceived
	ch <- collector.TotalPacketsSent
	ch <- collector.TotalPacketsReceived
}

func (collector *FritzBoxCollector) Collect(ch chan<- prometheus.Metric) {
	uPnPClient := newUPnPClient(viper.GetString("url"), viper.GetString("username"), viper.GetString("password"))
	metrics := uPnPClient.execute()

	extract := func(metric string) float64 {
		for k, v := range metrics {
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
	ch <- prometheus.MustNewConstMetric(collector.PhysicalLinkStatus, prometheus.GaugeValue, extract("PhysicalLinkStatus"))
	ch <- prometheus.MustNewConstMetric(collector.Layer1DownstreamMaxBitRate, prometheus.GaugeValue, extract("Layer1DownstreamMaxBitRate"))
	ch <- prometheus.MustNewConstMetric(collector.Layer1UpstreamMaxBitRate, prometheus.GaugeValue, extract("Layer1UpstreamMaxBitRate"))
	ch <- prometheus.MustNewConstMetric(collector.TotalBytesSent, prometheus.GaugeValue, extract("TotalBytesSent"))
	ch <- prometheus.MustNewConstMetric(collector.TotalBytesReceived, prometheus.GaugeValue, extract("TotalBytesReceived"))
	ch <- prometheus.MustNewConstMetric(collector.TotalPacketsSent, prometheus.GaugeValue, extract("TotalPacketsSent"))
	ch <- prometheus.MustNewConstMetric(collector.TotalPacketsReceived, prometheus.GaugeValue, extract("TotalPacketsReceived"))
}
