#!/usr/bin/python

# basically the sample code from the docs:
# https://pypi.python.org/pypi/rpiSht1x/1.2

from sht1x.Sht1x import Sht1x as SHT1x
dataPin = 16
clkPin = 7
sht1x = SHT1x(dataPin, clkPin, SHT1x.GPIO_BOARD)

temperature = sht1x.read_temperature_C()
humidity = sht1x.read_humidity()
dp = sht1x.calculate_dew_point(temperature, humidity)

print "Temp: %.2f C\tHumidity: %.2f %%\tDew Point: %.2f C" % (temperature, humidity, dp)
