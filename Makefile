image:
	docker build --no-cache -t fritzbox-upnp-exporter -t rekzi/fritzbox-upnp-exporter .
push:
	docker push rekzi/fritzbox-upnp-exporter:latest