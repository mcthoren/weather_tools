#!/bin/bash

# this script is run from cron

WT_DIR="/home/ghz/wx"
LOCK="$WT_DIR/plots/LOCK.wx"

[ -e $LOCK ] && {
	echo "$0: lock exists" | logger
	exit 1
}

# lock is also checked for and deleted on boot, in case of a crash
touch $LOCK

$WT_DIR/sensor_data_grab.py

$WT_DIR/grab_48h /home/ghz/wx/data bme680.dat
$WT_DIR/grab_48h /home/ghz/wx/data derived.dat
$WT_DIR/grab_48h /home/ghz/wx/data pi_temp

cd /home/ghz/wx/plots || exit 1
gnuplot $WT_DIR/cutie.wx.gnuplot

sync

/usr/bin/rsync -e "ssh -q" --timeout=60 -ur $WT_DIR/* wx1_sync:/wx1/ 2>/dev/null

rm $LOCK
