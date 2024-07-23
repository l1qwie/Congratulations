build:
	docker build . -t employees

remove:
	docker compose stop \
	&& docker compose rm \
	&& sudo rm -rf kafkadata/

up-debug:
	REDIS_LOGLEVEL=debug docker compose -f docker-compose.yml up --force-recreate

up-verbose:
	REDIS_LOGLEVEL=verbose docker compose -f docker-compose.yml up --force-recreate

up-notice:
	REDIS_LOGLEVEL=notice docker compose -f docker-compose.yml up --force-recreate

up-warning:
	REDIS_LOGLEVEL=warning docker compose -f docker-compose.yml up --force-recreate
