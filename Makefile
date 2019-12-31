image:
	docker build --no-cache -t fritzbox-upnp-exporter -t rekzi/fritzbox-upnp-exporter -f docker/Dockerfile .
push:
	docker push rekzi/fritzbox-upnp-exporter:latest
up:
	docker-compose -f docker/docker-compose.yml up
down:
	docker-compose -f docker/docker-compose.yml down

