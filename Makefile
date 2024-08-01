net:
	docker network create congratulations
launch tools:
	docker compose -f docker-compose-tools.yml up -d
delete:
	docker compose -f docker-compose-tools.yml stop \
	&& docker compose -f docker-compose-tools.yml rm -f \
	&& sudo rm -rf pgdata/
