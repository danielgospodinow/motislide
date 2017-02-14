#!/bin/bash

crontab -l > currentCron

sed --in-place '/changeWallpaper.sh/d' currentCron

crontab currentCron
rm currentCron

echo "Slideshow schedules removed successfully!"
