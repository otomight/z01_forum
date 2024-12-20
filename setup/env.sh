#!/bin/bash

function create_env() {
	ENV_EXAMPLE_FILE=$1
	ENV_FILE=$2

	if [ ! -f "$ENV_FILE" ]; then
		cat $ENV_EXAMPLE_FILE > $ENV_FILE
		echo -e "\e[33mWARNING: $ENV_FILE file created. Pls fill it with the right values!\e[0m"
	fi
}

create_env "$1" "$2"
