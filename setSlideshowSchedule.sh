#!/bin/bash

#Enter minutes for first argument. (We assume that the normal transition time is from 1 to 60 minutes)

requered_args=1
minutes=$1

if [ $# -lt $requered_args ]; then
	echo "Pass a single argument representing the transition time!"
	echo "------> ${0} <minutes>"
	exit 1
fi

crontab -l > currentCron

echo "*/${minutes} * * * * ${PWD}/changeWallpaper.sh ${PWD}/images" >> currentCron

crontab currentCron
rm currentCron

echo "Scheduler set successfully, enjoy!"
