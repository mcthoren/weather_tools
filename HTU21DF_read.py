#!/usr/bin/python

# written with help from:
# https://github.com/adafruit/Adafruit_HTU21DF_Library/
# https://github.com/adafruit/Adafruit_Python_BMP/
# and
# https://github.com/dalexgray/RaspberryPI_HTU21DF/

import Adafruit_GPIO.I2C as I2C
import time
import datetime
import math

temp_read = 0xE3
hum_read = 0xE5
write_reg = 0xE6
soft_reset = 0xFE

# constants from the datasheet
const_a = 8.1332
const_b = 1762.39
const_c = 235.66

htu_addy = 0x40
bus = 0

array = [0x00, 0x00]

dev0 = I2C.Device(htu_addy, bus)
dev0.write8(soft_reset, write_reg) # does this actually reset it? prly not....
time.sleep(0.2)

ts =  datetime.datetime.fromtimestamp(time.time()).strftime("%Y%m%d%H%M%S")

array = dev0.readList(temp_read, 2)
t0 = (array[0] * 256.0) + array[1]
temp = ((t0 * 175.72) / 65536) - 46.85

array = dev0.readList(hum_read, 2)
h0 = (array[0] * 256.0) + array[1]
hum = ((h0 * 125) / 65536) - 6

P_part = 10 ** (const_a - (const_b / (temp + const_c)))
T_dew = -1 * ((const_b / (math.log10(hum * (P_part / 100)) - const_a)) + const_c)

print "%s\tTemp: %.2f C\tHumidity: %.2f %%\tDew Point: %.2f C" % (ts, temp, hum, T_dew)
