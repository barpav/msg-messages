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
# make message KEY=session-key TO=userId TXT="Message text"
message:
	curl -v -X POST	-H "Content-Type: application/vnd.newPersonalMessage.v1+json" \
	-H "Authorization: Bearer $(KEY)" \
	-d '{"to": "$(TO)", "text": "$(TXT)"}' \
	localhost:8080
# make message-f KEY=session-key TO=userId TXT="Message text" F=file-id
message-f:
	curl -v -X POST	-H "Content-Type: application/vnd.newPersonalMessage.v1+json" \
	-H "Authorization: Bearer $(KEY)" \
	-d '{"to": "$(TO)", "text": "$(TXT)", "files": ["$(F)"]}' \
	localhost:8080
# make sync KEY=session-key A=after L=limit
sync:
	curl -v -H "Authorization: Bearer $(KEY)" \
	-H "Accept: application/vnd.messageUpdates.v1+json" \
	"localhost:8080?after=$(A)&limit=$(L)"
# make get-message KEY=session-key ID=message-id
get-message:
	curl -v -H "Authorization: Bearer $(KEY)" \
	"localhost:8080/$(ID)"