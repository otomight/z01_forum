#!/bin/bash

ACTION=$1
SASS_COMMAND=$2

PID_FILE="setup/.pid"

get_sass_watcher() {
	if [ -f "$PID_FILE" ]; then
		WATCHER_PID=$(cat "$PID_FILE")
		if ps -p "$WATCHER_PID" > /dev/null 2>&1; then
			echo "$WATCHER_PID"
			return 0
		fi
	fi
	return 1
}

start_sass_watcher() {
	if get_sass_watcher; then
		echo "Sass watcher still running with PID: $(get_sass_watcher)"
	else
		$SASS_COMMAND &
		WATCHER_PID=$!
		echo "$WATCHER_PID" > "$PID_FILE"
		echo "Sass watcher started with PID: $WATCHER_PID"
	fi
}

stop_sass_watcher() {
	if get_sass_watcher; then
		WATCHER_PID=$(get_sass_watcher)
		kill -9 "$WATCHER_PID" && rm -f "$PID_FILE"
		echo "Sass watcher process stopped with PID: $WATCHER_PID"
	else
		echo "No Sass watcher process found."
	fi
}

if [ "$ACTION" = "start" ]; then
	start_sass_watcher
elif [ "$ACTION" = "stop" ]; then
	stop_sass_watcher
fi
