#!/bin/sh

# meant to be called from cron every minute or so

/import/home/ghz/repos/weather_tools/sensor_data_grab.py

/import/home/ghz/repos/weather_tools/grab_48h /home/ghz/wx/data bme680.dat
/import/home/ghz/repos/weather_tools/grab_48h /home/ghz/wx/data derived.dat
/import/home/ghz/repos/weather_tools/grab_48h /home/ghz/wx/data pi_temp
