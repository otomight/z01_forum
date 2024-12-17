IMAGE_NAME=forum_i
CONTAINER_NAME=forum_c
PORT=443

all: build run

audit: clean build images run

containers:
	sudo docker ps -a

images:
	sudo docker images

build:
	sudo docker build -t $(IMAGE_NAME) .

run:
	sudo docker run -p $(PORT):$(PORT) --name $(CONTAINER_NAME) $(IMAGE_NAME)

stop:
	sudo docker stop $(CONTAINER_NAME)

rm:
	sudo docker rm $(CONTAINER_NAME)

rmi:
	sudo docker rmi $(IMAGE_NAME)

clean: stop rm rmi

re: clean all

.PHONY: all build run stop rm rmi clean re
