OS_NAME	:= $(if $(filter Windows_NT,$(OS)),WINDOWS,UNIX)

#	UNIX
RM_CMD_UNIX				:= rm -f
BIN_FILE_UNIX			:= main
ENV_SETUP_FILE_UNIX		= $(SETUP_DIR)/env.sh
SASS_WATCHER_FILE_UNIX	= $(SETUP_DIR)/watcher.sh

#	WINDOWS
RM_CMD_WINDOWS				:= del /Q
BIN_FILE_WINDOWS			:= main.exe
ENV_SETUP_FILE_WINDOWS		= $(SETUP_DIR)/env.ps1
SASS_WATCHER_FILE_WINDOWS	= $(SETUP_DIR)/watcher.ps1



#	VARIABLES
IMAGE_NAME		:= forum_i
CONTAINER_NAME	:= forum_c
PORT			:= 443


#	DIRECTORIES
SETUP_DIR			:= setup


#	FILES
POSTS_IMAGES_FILES		:= data/images/posts/*
DB_FILE					:= forum.db
BIN_FILE				:= $(BIN_FILE_$(OS_NAME))

# Setup
ENV_FILE				:= .env
ENV_EXAMPLE_FILE		:= $(SETUP_DIR)/.env.example
ENV_SETUP_FILE			:= $(ENV_SETUP_FILE_$(OS_NAME))
SASS_WATCHER_FILE		:= $(SASS_WATCHER_FILE_$(OS_NAME))

# HTTPS configs
OPEN_SSL_CNF_FILE		:= $(SETUP_DIR)/openssl.cnf
CERTOUT_FILE			:= server.crt
KEYOUT_FILE				:= server.key
MAIN_SCSS_FILE			:= web/src/scss/main.scss
MAIN_CSS_OUT_FILE		:= web/static/style/main.css


#	KEYS
START_SASS_WATCHER_KEY	:= start
STOP_SASS_WATCHER_KEY	:= stop

#	COMMANDS
SASS_CMD	:= npx sass --watch $(MAIN_SCSS_FILE):$(MAIN_CSS_OUT_FILE) \
				--style=compressed
RM_CMD		= $(RM_CMD_$(OS_NAME))


#	BASIC RULES
all: gencertif build run

re:	clean clean-data all

build: create-env-file compile-ts compile-scss
	@echo Build programm binary...
	@go build -o $(BIN_FILE)

run:
	@echo Run application...
ifeq ($(OS_NAME),WINDOWS)
	@powershell -Command "Start-Process '.\\$(BIN_FILE)' -Verb RunAs"
else
	@sudo ./$(BIN_FILE)
endif

clean: sass-watch-stop
	$(RM_CMD) $(BIN_FILE)

clean-data:
	$(RM_CMD) $(DB_FILE)
ifeq ($(OS_NAME),WINDOWS)
	@powershell Remove-Item -Path "$(POSTS_IMAGES_FILES)" -Force
else
	sudo $(RM_CMD) $(POSTS_IMAGES_FILES)
endif

fclean: clean clean-data
	$(RM_CMD) $(CERTOUT_FILE) $(KEYOUT_FILE)


#	SPECIAL RULES
gencertif:
ifeq ($(wildcard $(CERTOUT_FILE) $(KEYOUT_FILE)),)
	openssl req -x509 -config $(OPEN_SSL_CNF_FILE) \
			-out $(CERTOUT_FILE) -keyout $(KEYOUT_FILE)
endif

create-env-file:
ifeq ($(OS_NAME),WINDOWS)
	@powershell -ExecutionPolicy Bypass -File $(ENV_SETUP_FILE) \
				-EnvExampleFile $(ENV_EXAMPLE_FILE) \
				-EnvFile $(ENV_FILE)
else
	@./$(ENV_SETUP_FILE) $(ENV_EXAMPLE_FILE) $(ENV_FILE)
endif

compile-ts:
	@echo Compiling typescript to javascript...
	@npx tsc

compile-scss:
	@echo Compiling scss to css and start edit watcher...
ifeq ($(OS_NAME),WINDOWS)
	@powershell -ExecutionPolicy Bypass -File $(SASS_WATCHER_FILE) \
				-Action "$(START_SASS_WATCHER_KEY)" \
				-SassCommand "$(SASS_CMD)"
else
	@./$(SASS_WATCHER_FILE) "$(START_SASS_WATCHER_KEY)" "$(SASS_CMD)"
endif

sass-watch-stop:
ifeq ($(OS_NAME),WINDOWS)
	@powershell -ExecutionPolicy Bypass -File $(SASS_WATCHER_FILE) \
				-Action "$(STOP_SASS_WATCHER_KEY)"
else
	@./$(SASS_WATCHER_FILE) "$(STOP_SASS_WATCHER_KEY)" "$(SASS_CMD)"
endif


#	DOCKER RULES
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

.PHONY: all gencertif build run watch clean re \
		dall dbuild drun dstop drm drmi dclean dre
