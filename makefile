DOCKER_COMPOSE_COMMAND = docker compose -f ./build/docker-compose.yaml -p="zylo"

redis: 
	$(DOCKER_COMPOSE_COMMAND) up -d redis

server: 
	$(DOCKER_COMPOSE_COMMAND) up -d api