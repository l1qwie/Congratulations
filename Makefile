net:
	docker network create congratulations
rm:
	docker compose stop \
	&& docker compose rm \
	&& sudo rm -rf pgdata/
up-all:
	docker compose -f docker-compose.yml up --force-recreate
up-empl:
	docker compose -f docker-compose-employees.yml up --force-recreate
rb-empl:
	docker build . -t congratulations-empl