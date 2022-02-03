#! /bin/sh

docker-compose build --pull
docker-compose run --rm run-example
