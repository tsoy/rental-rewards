# Variables
PUBSUB_EMULATOR_HOST ?= localhost:8085

.PHONY: help pubsub-init

help:
	@echo "Available commands:\n"
	@grep -E '^##' Makefile | sed -e 's/## //'

## make run          : air .
run:
	air .

## make docker-up    : Start the DB and Pub/Sub emulator in Docker
dup:
	docker compose -f docker-compose-dev.yml up -d

## make docker-down  : Stop the DB and Pub/Sub emulator in Docker
ddown:
	docker compose -f docker-compose-dev.yml down

## make pubsub-init  : Initialize topics & subscriptions in Pub/Sub emulator
pubsub-init:
	curl -X PUT "http://$(PUBSUB_EMULATOR_HOST)/v1/projects/test-project/topics/payment.completed"
	curl -X PUT "http://$(PUBSUB_EMULATOR_HOST)/v1/projects/test-project/subscriptions/rewards-worker-sub" \
	  -H "Content-Type: application/json" \
	  -d '{"topic": "projects/test-project/topics/payment.completed"}'

## Execute migrations
mig:
	migrate -path=./migrations -database=$RR_DB_DSN up

#docker container exec -it rental-rewards-db-1 psql