#!/bin/bash

# meant to be called from cron every minute or so

LOCK="/home/ghz/wx/wx.lock"

[ -e "${LOCK}" ] && {
	echo "$0: lock exists" | logger
	exit 1
}

# lock is also checked for and deleted on boot, in case of a crash
touch "${LOCK}"

HOST_N="$(hostname -s)"
[[ "${HOST_N}" == "keen" || "${HOST_N}" == "infinity" ]] && WT_DIR='/import/home/ghz/repos/weather_tools/'
[[ "${HOST_N}" == "cutie" ]] && WT_DIR='/home/ghz/wx'

$WT_DIR/sensor_data_grab.py

$WT_DIR/grab_48h /home/ghz/wx/data bme680.dat
$WT_DIR/grab_48h /home/ghz/wx/data derived.dat
$WT_DIR/grab_48h /home/ghz/wx/data pi_temp

cd /home/ghz/wx/plots || exit 1
gnuplot "$WT_DIR/wx.gnuplot"

# this prevents us from losing 20m of data when the power drops.
sync

[[ "${HOST_N}" == "keen" ]] && {
	/usr/bin/rsync -ur --timeout=50 /home/ghz/wx /import/home/ghz/repos/dust_wx wx0_sync:/wx0/ 2> /dev/null
}

[[ "${HOST_N}" == "cutie" ]] && {
	/usr/bin/rsync -e "ssh -q" --timeout=60 -ur $WT_DIR/* wx1_sync:/wx1/ 2>/dev/null
}

[[ "${HOST_N}" == "infinity" ]] && {
	/usr/bin/rsync -ur --delete --timeout=50 --exclude='.git' /home/ghz/wx/ \
	/import/home/ghz/repos/weather_tools/ wx_0x0b_sync:/wx_0x0b/ # 2> /dev/null
}

rm "${LOCK}"
