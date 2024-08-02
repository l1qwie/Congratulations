net:
	docker network create congratulations
launch-test:
	docker compose -f docker-compose-test.yml up -d
launch-all:
	docker compose -f docker-compose.yml up -d
delete:
	docker compose -f docker-compose.yml stop \
	&& docker compose -f docker-compose.yml rm -f \
	&& sudo rm -rf pgdata/
build:
	docker build . -t congratulations
up:
	docker run --rm --name notif --network congratulations congratulations /app/bin