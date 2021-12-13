#!/bin/sh

# meant to be called from cron every minute or so

LOCK="/home/ghz/wx/wx.lock"

[ -e "${LOCK}" ] && {
	echo "$0: lock exists" | logger
	exit 1
}

# lock is also checked for and deleted on boot, in case of a crash
touch "${LOCK}"

WT_DIR='/import/home/ghz/repos/weather_tools/'

$WT_DIR/sensor_data_grab.py

$WT_DIR/grab_48h /home/ghz/wx/data bme680.dat
$WT_DIR/grab_48h /home/ghz/wx/data derived.dat
$WT_DIR/grab_48h /home/ghz/wx/data pi_temp

cd /home/ghz/wx/plots || exit 1
gnuplot "$WT_DIR/keen.wx.gnuplot"

sync

/usr/bin/rsync -ur --delete --timeout=50 /home/ghz/wx /import/home/ghz/repos/dust_wx wx0_sync:/wx0/ 2> /dev/null

rm "${LOCK}"
