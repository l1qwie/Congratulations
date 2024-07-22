net:
	docker network create congratulations
rm:
	docker compose stop \
	&& docker compose rm \
	&& sudo rm -rf pgdata/
up-all:
	docker compose -f docker-compose.yml up --force-recreate
up-empl-debug:
	REDIS_LOGLEVEL=debug docker compose -f docker-compose-employees.yml up --force-recreate
up-empl-verbose:
	REDIS_LOGLEVEL=verbose docker compose -f docker-compose-employees.yml up --force-recreate
up-empl-notice:
	REDIS_LOGLEVEL=notice docker compose -f docker-compose-employees.yml up --force-recreate
up-empl-warning:
	REDIS_LOGLEVEL=warning docker compose -f docker-compose-employees.yml up --force-recreate
rb-empl:
	docker build . -t congratulations-empl