# If you are using Visual Studio Code for development,
# you can start by re-launching this folder as a Docker container.
# All necessary configurations are already provided in the [.devcontainer] directory. 
# If you are new to this, you can follow this video tutorial to get started: 
# [youtube video link here]

# OR

# If you are using any other IDE the following commands may help you

server_image_name := readme-studio-image
dev_image_name := readmestudiodockercompose-dev

container_name := readme-studio-container

ls_dangling_cmd := docker images -f "dangling=true"

dev:
	docker compose -f dev.docker-compose.yaml up dev --detach

devsh:
	docker exec -it ${container_name} /bin/sh

server:
	docker build -t ${server_image_name} -f server.Dockerfile .
	docker run --name ${container_name} ${server_image_name}

dangling_list:
	${ls_dangling_cmd}

# This command will clean most images & containers
clean: 
	-docker compose -f dev.docker-compose.yaml down
	-docker stop ${container_name}
	-docker rm ${container_name}
	-docker rmi ${dev_image_name}
	-docker rmi ${server_image_name}
	-docker rmi `${ls_dangling_cmd} -q`
