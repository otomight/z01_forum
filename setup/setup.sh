#!/bin/bash

function create_env() {
	ENV_EXAMPLE_FILE=$1
	ENV_FILE=$2

	if [ ! -f "$ENV_FILE" ]; then
		cat $ENV_EXAMPLE_FILE > $ENV_FILE
		echo "$ENV_FILE file created. Pls fill it with the right values"
	fi
}

create_env "$1" "$2"
