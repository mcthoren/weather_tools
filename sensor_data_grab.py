#!/usr/bin/python
# -*- coding: utf-8 -*-

import time
import datetime

wx_dir = "/home/ghz/wx"

def write_out(file_name, data, mode):
	out_file_fd = open(file_name, mode)
	out_file_fd.write(data)
	out_file_fd.close()

def write_out_dat_stamp(ts, n_plate, data):
	# year directories should be created once a year from cron
	# that way we aren't unnecessarily checking for one every minute of every day for a year

	f_ts = ts[0:8]
	y_ts = ts[0:4]
	write_out(wx_dir+'/data/'+y_ts+'/'+n_plate+'.'+f_ts, data, 'a')

def bme680_read():
	# sensor + breakout board from:
	# https://www.adafruit.com/product/3660

	# libraries and examples from:
	# https://learn.pimoroni.com/tutorial/sandyj/getting-started-with-bme680-breakout

	# dew point equations and constants from:
	# http://journals.ametsoc.org/doi/pdf/10.1175/BAMS-86-2-225

	import bme680
	import math

	pressures = []
	iter = 16
	avg = 0

	Ca = 17.625
	Cb = 243.04 # [°C]

	sensor = bme680.BME680(i2c_addr=0x77)

	sensor.set_humidity_oversample(bme680.OS_2X)
	sensor.set_pressure_oversample(bme680.OS_4X)
	sensor.set_temperature_oversample(bme680.OS_8X)
	# sensor.set_filter(bme680.FILTER_SIZE_3)
	sensor.set_gas_status(bme680.ENABLE_GAS_MEAS)

	sensor.set_gas_heater_temperature(320)
	sensor.set_gas_heater_duration(150)
	sensor.select_gas_heater_profile(0)

	ts = datetime.datetime.fromtimestamp(time.time()).strftime("%Y%m%d%H%M%S")

	temp = sensor.data.temperature
	hum = sensor.data.humidity

	gamma = math.log(hum / 100) + ((Ca * temp) / (Cb + temp))
	Tdew = (Cb * gamma) / (Ca - gamma)

	# go back to averaging pressure samples as in the bmp085 datasheet
	for x in range(0, iter):
		pressures.append(sensor.data.pressure)
		avg += pressures[x]

	pres_avg = avg / iter # [hPa]

	gas_res = 0
	tries = 5

	while (tries > 0):
		sensor.get_sensor_data()
		if sensor.data.heat_stable:
			gas_res = sensor.data.gas_resistance
		else:
			time.sleep(1)

		tries = tries - 1

	dat_string = "%s\tTemp: %.2f C\tHumidity: %.2f %%\tPressure: %.3f kPa\tAirQ: %d Ohms\tTdew: %.2f C\n" % (ts, temp, hum, pres_avg / 10, gas_res, Tdew)

	write_out_dat_stamp(ts, 'bme680.dat', dat_string)

	return (temp, hum, pres_avg / 10, gas_res, Tdew)
		
def pi_temp_read():
	temp_file = "/sys/class/thermal/thermal_zone0/temp"
	temp_file_fd = open(temp_file, 'r')

	ts = datetime.datetime.fromtimestamp(time.time()).strftime("%Y%m%d%H%M%S")

	temp_data = temp_file_fd.read()
	temp_file_fd.close()

	dat_string = "%s\t%s" % (ts, temp_data)

	write_out_dat_stamp(ts, 'pi_temp', dat_string)

	return (float(temp_data) / 1000)

def gen_index(etemp, ehum, press, pitemp, edp):
	plate = wx_dir+"/wx_index.html.template"
	plate_fd = open(plate, 'r')
	plate_dat = plate_fd.read()
	plate_fd.close()

	ts = datetime.datetime.fromtimestamp(time.time()).strftime("%FT%TZ")

	plate_dat = plate_dat.replace("EXTTEMP", str("%.2f" % etemp))
	plate_dat = plate_dat.replace("EXTHUM", str("%.2f" % ehum))
	plate_dat = plate_dat.replace("EXTDP", str("%.2f" % edp))
	plate_dat = plate_dat.replace("REL_PRESS", str("%.3f" % press))
	plate_dat = plate_dat.replace("PITEMP", str("%.2f" % pitemp))
	plate_dat = plate_dat.replace("DATE", ts)

	write_out(wx_dir+'/plots/wx.html', plate_dat, 'w')

if __name__ == "__main__":
	press_cal = 0.37 # kPa
	pi_temp = pi_temp_read()
	(e_temp, e_hum, press, gas_r, Tdew) = bme680_read()
	gen_index(e_temp, e_hum, press + press_cal, pi_temp, Tdew)
