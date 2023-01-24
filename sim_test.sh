#!/bin/bash

kill_process() {
	PID_GO=$(ps aufx | grep main.go | grep -v color |grep -v grep| awk -F' ' '{ print $2 }')
	if [ -n "${PID_GO}" ]; then
		kill ${PID_GO}
		wait ${PID_GO} 2>/dev/null
	fi
	
	PID_PYTHON=$(ps aufx | grep app.py | grep -v color |grep -v grep| awk -F' ' '{ print $2 }')
	if [ -n "${PID_PYTHON}" ]; then
		kill ${PID_PYTHON}
		wait ${PID_PYTHON} 2>/dev/null
	fi
}

main() {
	kill_process
	echo `date`

	HOURS=24
	POLLING_INTERVAL=5
	POLLING_INTERVAL_SECS=$((5*60))
	TIME=0

	TIME_TICKER=$(($HOURS*60/$POLLING_INTERVAL))

	PRESENT_WD=$(pwd)

	cd ${PRESENT_WD}/simulator/src
	python3.9 app.py &
	sleep 3
	sudo timedatectl set-ntp 0
	sudo date -s "$(date +"%Y-%m-%d") 00:00:00 0 seconds" +"%H:%M:%S"
	cd ${PRESENT_WD}
	go run main.go &
	for i in $(seq 1 $TIME_TICKER); do
		sleep 1
		sudo date -s "$(date +"%Y-%m-%d") 00:00:00 $INTERVAL seconds" +"%H:%M:%S"
		echo `date`
		INTERVAL=$(($i*$POLLING_INTERVAL_SECS))
	done
	sudo timedatectl set-ntp 1
	kill_process
	sleep 2
}

main
