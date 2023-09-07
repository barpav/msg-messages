build-m:
	sudo docker image rm -f ghcr.io/barpav/msg-messages:v1
	sudo docker build -t ghcr.io/barpav/msg-messages:v1 -f docker/service/Dockerfile .
	sudo docker image ls
build-s:
	sudo docker image rm -f ghcr.io/barpav/msg-storage-messages:v1
	sudo docker build -t ghcr.io/barpav/msg-storage-messages:v1 -f docker/storage/Dockerfile ./docker/storage
	sudo docker image ls

push-m:
	sudo docker push ghcr.io/barpav/msg-messages:v1
push-s:
	sudo docker push ghcr.io/barpav/msg-storage-messages:v1
push:
	sudo docker push ghcr.io/barpav/msg-messages:v1
	sudo docker push ghcr.io/barpav/msg-storage-messages:v1

up:
	sudo docker-compose up -d --wait
down:
	sudo docker-compose down

up-debug:
	sudo docker-compose -f compose-debug.yaml up -d --wait
down-debug:
	sudo docker-compose -f compose-debug.yaml down

user:
	curl -v -X POST	-H "Content-Type: application/vnd.newUser.v1+json" \
	-d '{"id": "jane", "name": "Jane Doe", "password": "My1stGoodPassword"}' \
	localhost:8081
session:
	curl -v -X POST -H "Authorization: Basic amFuZTpNeTFzdEdvb2RQYXNzd29yZA==" localhost:8082