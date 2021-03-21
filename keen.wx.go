#!/bin/sh

# meant to be called from cron every minute or so

WT_DIR='/import/home/ghz/repos/weather_tools/'

$WT_DIR/sensor_data_grab.py

$WT_DIR/grab_48h /home/ghz/wx/data bme680.dat
$WT_DIR/grab_48h /home/ghz/wx/data derived.dat
$WT_DIR/grab_48h /home/ghz/wx/data pi_temp

cd /home/ghz/wx/plots || exit 1
gnuplot "$WT_DIR/keen.wx.gnuplot"

sync
