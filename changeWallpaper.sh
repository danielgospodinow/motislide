#!/bin/sh

user=$(whoami)

fl=$(find /proc -maxdepth 2 -user $user -name environ -print -quit)
for i in {1..5}
do
  fl=$(find /proc -maxdepth 2 -user $user -name environ -newer "$fl" -print -quit)
done

export DBUS_SESSION_BUS_ADDRESS=$(grep -z DBUS_SESSION_BUS_ADDRESS "$fl" | cut -d= -f2-)

setWallpaper()
{
	dconf write "/org/gnome/desktop/background/picture-uri" "'file://${1}'"
}

wallpapers=(${1}/*)
setWallpaper ${wallpapers[RANDOM % ${#wallpapers[@]}]}
