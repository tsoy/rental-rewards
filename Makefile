#docker container exec -it rental-rewards-db-1 psql
#migrate -path=./migrations -database=$RR_DB_DSN up