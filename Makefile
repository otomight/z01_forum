IMAGE_NAME=forum_i
CONTAINER_NAME=forum_c
PORT=443

CERTOUT_FILE=server.crt
KEYOUT_FILE=server.key
MAIN_SCSS_FILE=web/src/scss/main.scss
MAIN_CSS_OUT_FILE=web/static/style/main.css
DB_FILE=forum.db

ifeq ($(OS),Windows_NT)
	BIN_FILE=main.exe
	RM_CMD=del
else
	BIN_FILE=main
	RM_CMD=rm -f
endif

all: gencertif build run

re:	clean all

gencertif:
ifeq ($(wildcard $(CERTOUT_FILE) $(KEYOUT_FILE)),)
	openssl req -x509 -config openssl.cnf \
			-out $(CERTOUT_FILE) -keyout $(KEYOUT_FILE)
endif

build:
	npx tsc
	npx sass $(MAIN_SCSS_FILE):$(MAIN_CSS_OUT_FILE) --style=compressed
	go build -o $(BIN_FILE)

run:
ifeq ($(OS),Windows_NT)
	powershell -Command "Start-Process '.\\$(BIN_FILE)' -Verb RunAs"
else
	sudo ./$(BIN_FILE)
endif

clean:
	$(RM_CMD) $(BIN_FILE)

fclean: clean
	$(RM_CMD) $(CERTOUT_FILE) $(KEYOUT_FILE) $(DB_FILE)

# DOCKER
dall: dbuild drun

dcontainers:
	sudo docker ps -a

dimages:
	sudo docker images

dbuild:
	sudo docker build \
	--build-arg CERTOUT_FILE=$(CERTOUT_FILE) \
	--build-arg KEYOUT_FILE=$(KEYOUT_FILE) \
	--build-arg MAIN_SCSS_FILE=$(MAIN_SCSS_FILE) \
	--build-arg MAIN_CSS_OUT_FILE=$(MAIN_CSS_OUT_FILE) \
	-t $(IMAGE_NAME) .

drun:
	sudo docker run -p $(PORT):$(PORT) --name $(CONTAINER_NAME) $(IMAGE_NAME)

dstop:
	sudo docker stop $(CONTAINER_NAME)

drm:
	sudo docker rm $(CONTAINER_NAME)

drmi:
	sudo docker rmi $(IMAGE_NAME)

dclean: dstop drm drmi

dre: dclean dall

.PHONY: all gencertif build run clean re \
		dall dbuild drun dstop drm drmi dclean dre
