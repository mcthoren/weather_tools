#!/usr/bin/python
# -*- coding: utf-8 -*-

import time, datetime, sys

sys.path.append('/home/ghz/wxlib')
import wxlib as wx

wx_dir = "/home/ghz/wx"

def bme680_read():
	# sensor + breakout board from:
	# https://www.adafruit.com/product/3660

	# libraries and examples from:
	# https://learn.pimoroni.com/tutorial/sandyj/getting-started-with-bme680-breakout

	import bme680
	import math

	iter = 16
	avg = 0

	sensor = bme680.BME680(i2c_addr=0x77)

	sensor.set_humidity_oversample(bme680.OS_2X)
	sensor.set_pressure_oversample(bme680.OS_4X)
	sensor.set_temperature_oversample(bme680.OS_8X)
	# sensor.set_filter(bme680.FILTER_SIZE_3)
	sensor.set_gas_status(bme680.ENABLE_GAS_MEAS)

	sensor.set_gas_heater_temperature(320)
	sensor.set_gas_heater_duration(150)
	sensor.select_gas_heater_profile(0)

	temp = sensor.data.temperature
	hum = sensor.data.humidity

	# go back to averaging pressure samples as in the bmp085 datasheet
	for x in range(0, iter):
		avg += sensor.data.pressure

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

	return (temp, hum, pres_avg / 10, gas_res)
		
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

	wx.write_out(wx_dir+'/plots/wx.html', plate_dat, 'w')

if __name__ == "__main__":
	ts = datetime.datetime.fromtimestamp(time.time()).strftime("%Y%m%d%H%M%S")
	pi_temp = wx.pi_temp_read()
	pi_dat_string = "%s\t%s" % (ts, pi_temp)
	wx.write_out_dat_stamp(ts, 'pi_temp', pi_dat_string, wx_dir)

	press_cal = 5.900 # kPa
	(e_temp, e_hum, press, gas_r) = bme680_read()
	Tdew = wx.dew_point_c(e_temp, e_hum)

	bme_dat_string = \
	"%s\tTemp: %.2f C\tHumidity: %.2f %%\tPressure: %.3f kPa\tAirQ: %d Ohms\tTdew: %.2f C\n" \
	% (ts, e_temp, e_hum, press, gas_r, Tdew)

	wx.write_out_dat_stamp(ts, 'bme680.dat', bme_dat_string, wx_dir)

	abs_hum = wx.abs_hum_g_mmm(e_temp, e_hum)
	heat_i = wx.heat_index(e_temp, e_hum)
	if (heat_i == -1):
		heat_i = e_temp  # dirty hack cuz the model is pretty wonky outside a narrow range

	derived_dat_string = \
	"%s\tAbsolute Humidity: %.2f g/m³\tHeat Index: %.2f °C\tDew Point: %.2f °C\n" \
	% (ts, abs_hum, heat_i, Tdew)
	wx.write_out_dat_stamp(ts, 'derived.dat', derived_dat_string, wx_dir)

	gen_index(e_temp, e_hum, press + press_cal, float(pi_temp) / 1000, Tdew)
