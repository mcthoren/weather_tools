#!/usr/bin/python

# written with help from:
# https://github.com/adafruit/Adafruit_HTU21DF_Library/
# https://github.com/adafruit/Adafruit_Python_BMP/
# https://github.com/dalexgray/RaspberryPI_HTU21DF/

import time
import datetime

wx_dir = "/home/ghz/wx"

def htu21df_read():
	# https://www.adafruit.com/product/1899
	
	import Adafruit_GPIO.I2C as I2C
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
	f_ts =  datetime.datetime.fromtimestamp(time.time()).strftime("%Y%m%d")

	array = dev0.readList(temp_read, 2)
	t0 = (array[0] * 256.0) + array[1]
	temp = ((t0 * 175.72) / 65536) - 46.85

	array = dev0.readList(hum_read, 2)
	h0 = (array[0] * 256.0) + array[1]
	hum = ((h0 * 125) / 65536) - 6

	P_part = 10 ** (const_a - (const_b / (temp + const_c)))
	T_dew = -1 * ((const_b / (math.log10(hum * (P_part / 100)) - const_a)) + const_c)

	dat_string = "%s\tTemp: %.2f C\tHumidity: %.2f %%\tDew Point: %.2f C\n" % (ts, temp, hum, T_dew)

	out_file_n = wx_dir+'/data/htu21df_grab.dat.'+f_ts
	out_file_fd = open(out_file_n, 'a')
	out_file_fd.write(dat_string)
	out_file_fd.close()
	
	# print "%s\tTemp: %.2f C\tHumidity: %.2f %%\tDew Point: %.2f C" % (ts, temp, hum, T_dew)

	return (temp, hum, T_dew)

def bmp085_read():
	# https://www.adafruit.com/products/391
	import Adafruit_BMP.BMP085 as BMP085

	pressures = []
	iter = 16       
	avg = 0

	sensor = BMP085.BMP085(mode=BMP085.BMP085_ULTRAHIGHRES)

	ts = datetime.datetime.fromtimestamp(time.time()).strftime("%Y%m%d%H%M%S")
	f_ts =  datetime.datetime.fromtimestamp(time.time()).strftime("%Y%m%d")

	temp = sensor.read_temperature()

	# datasheet suggests averaging
	for x in range(0, iter):
		pressures.append(sensor.read_pressure())
		avg += pressures[x]

	avg = avg / iter

	# means sea level pressure calibration
	# calibrated against:
	# http://www.geomar.de/service/wetter/
	# [kPa]
	mslp_calibration = .749

	dat_string = "%s\t%.2f C\t\t%.3f kPa\n" % (ts, temp, ((avg / 1000.0) + mslp_calibration))

	out_file_n = wx_dir+'/data/bmp0085_grab.dat.'+f_ts
	out_file_fd = open(out_file_n, 'a')
	out_file_fd.write(dat_string)
	out_file_fd.close()

	# print "%s\t%.2f C\t\t%.3f kPa" % (ts, temp, ((avg / 1000.0) + mslp_calibration))

	return (temp, (avg / 1000.0) + mslp_calibration)

def sht11_read():
	# basically the sample code from the docs:
	# https://pypi.python.org/pypi/rpiSht1x/1.2

	from sht1x.Sht1x import Sht1x as SHT1x
	dataPin = 16
	clkPin = 7

	sht1x = SHT1x(dataPin, clkPin, SHT1x.GPIO_BOARD)

	ts = datetime.datetime.fromtimestamp(time.time()).strftime("%Y%m%d%H%M%S")
	f_ts =  datetime.datetime.fromtimestamp(time.time()).strftime("%Y%m%d")

	temperature = sht1x.read_temperature_C()
	humidity = sht1x.read_humidity()
	dp = sht1x.calculate_dew_point(temperature, humidity)

	dat_string = "%s\tTemp: %.2f C\tHumidity: %.2f %%\tDew Point: %.2f C\n" % (ts, temperature, humidity, dp)

	out_file_n = wx_dir+'/data/sht11_grab.dat.'+f_ts
	out_file_fd = open(out_file_n, 'a')
	out_file_fd.write(dat_string)
	out_file_fd.close()

	# print "Temp: %.2f C\tHumidity: %.2f %%\tDew Point: %.2f C" % (temperature, humidity, dp)

	return (temperature, humidity, dp)

def pi_temp_read():
	temp_file = "/sys/class/thermal/thermal_zone0/temp"
	temp_file_fd = open(temp_file, 'r')

	ts = datetime.datetime.fromtimestamp(time.time()).strftime("%Y%m%d%H%M%S")
	f_ts =  datetime.datetime.fromtimestamp(time.time()).strftime("%Y%m%d")

	temp_data = temp_file_fd.read()
	temp_file_fd.close()

	dat_string = "%s\t%s" % (ts, temp_data)

	out_file_n = wx_dir+'/data/pi_temp.'+f_ts
	out_file_fd = open(out_file_n, 'a')
	out_file_fd.write(dat_string)
	out_file_fd.close()

	return (float(temp_data) / 1000)

def gen_index(etemp, ehum, edp, press, bmptemp, itemp, ihum, idp, pitemp):
	plate = wx_dir+"/wx_index.html.template"
	plate_fd = open(plate, 'r')
	plate_dat = plate_fd.read()
	plate_fd.close()

	ts = datetime.datetime.fromtimestamp(time.time()).strftime("%FT%TZ")

	plate_dat = plate_dat.replace("EXTTEMP", str("%.2f" % etemp))
	plate_dat = plate_dat.replace("EXTHUM", str("%.2f" % ehum))
	plate_dat = plate_dat.replace("EXTDP", str("%.2f" % edp))
	plate_dat = plate_dat.replace("REL_PRESS", str("%.3f" % press))
	plate_dat = plate_dat.replace("INTTEMPSHT", str("%.2f" % itemp))
	plate_dat = plate_dat.replace("INTHUMSHT", str("%.2f" % ihum))
	plate_dat = plate_dat.replace("INTTEMP", str("%.2f" % bmptemp))
	plate_dat = plate_dat.replace("PITEMP", str("%.2f" % pitemp))
	plate_dat = plate_dat.replace("DATE", ts)

	out_file_n = wx_dir+'/plots/wx.html'
	out_file_fd = open(out_file_n, 'w')
	out_file_fd.write(plate_dat)
	out_file_fd.close()

if __name__ == "__main__":
	(bmp_temp, press) = bmp085_read()
	(e_temp, e_hum, e_dp) = htu21df_read()
	(i_temp, i_hum, i_dp) = sht11_read()
	pi_temp = pi_temp_read()
	gen_index(e_temp, e_hum, e_dp, press, bmp_temp, i_temp, i_hum, i_dp, pi_temp)
