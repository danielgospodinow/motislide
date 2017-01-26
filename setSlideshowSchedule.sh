#!/bin/bash

#Enter minutes for first argument. (We assume that the normal transition time is from 1 to 60 minutes)

minutes=$1

crontab -l > currentCron

echo "*/${minutes} * * * * ${PWD}/changeWallpaper.sh ${PWD}/images" >> currentCron

crontab currentCron
rm currentCron
