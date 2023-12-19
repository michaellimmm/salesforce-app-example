setup:
	go mod download

gen-proto:
	buf generate proto

run-dev-docker-mongo:
	docker-compose -f ./docker-compose.yaml up -d mongo --remove-orphans