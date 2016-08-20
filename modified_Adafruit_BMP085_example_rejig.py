#!/usr/bin/python

# Example Code from Adafruit's BMP085 examples.
# https://www.adafruit.com/products/391

# my hacked up version of Adafruit's sample code.
# no, i haven't ever written anything in python before.

import Adafruit_BMP.BMP085 as BMP085
import time
import datetime

pressures = []
iter = 16       
avg = 0

sensor = BMP085.BMP085(mode=BMP085.BMP085_ULTRAHIGHRES)

ts = datetime.datetime.fromtimestamp(time.time()).strftime("%Y%m%d%H%M%S")
temp = sensor.read_temperature()

# ideally i wouldn't average, i would look for successive equal values,
# but this sensor doesn't seem very amenable to that idea.
# turns out datasheet suggests averaging
for x in range(0, iter):
	pressures.append(sensor.read_pressure())
	avg += pressures[x]

avg = avg / iter

# means sea level pressure calibration
# calibrated against:
# http://www.geomar.de/service/wetter/
# kPa
# mslp_calibration = .832

# moved up a few floors
mslp_calibration = .749

print "%s\t%.2f C\t\t%.3f kPa" % (ts, temp, ((avg / 1000.0) + mslp_calibration))
