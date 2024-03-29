#!/usr/bin/python3
# -*- coding: utf-8 -*-

import time, sys, os

sys.path.append('/home/ghz/wxlib')
import wxlib as wx

wx_dir = "/home/ghz/wx"

def gen_index(etemp, ehum, press, pitemp, edp, abs_hum, heat_i, Tw, title):
	plate = wx_dir+"/wx_index.html.template"
	plate_fd = open(plate, 'r')
	plate_dat = plate_fd.read()
	plate_fd.close()

	ts = time.strftime("%FT%TZ", time.gmtime())

	plate_dat = plate_dat.replace("EXTTEMP", str("%.2f" % etemp))
	plate_dat = plate_dat.replace("EXTHUM", str("%.2f" % ehum))
	plate_dat = plate_dat.replace("EXTDP", str("%.2f" % edp))
	plate_dat = plate_dat.replace("REL_PRESS", str("%.3f" % press))
	plate_dat = plate_dat.replace("PITEMP", str("%.2f" % pitemp))
	plate_dat = plate_dat.replace("ABSHUM", str("%.2f" % abs_hum))
	plate_dat = plate_dat.replace("HTIDX", str("%.2f" % heat_i))
	plate_dat = plate_dat.replace("WBT", str("%.2f" % Tw))
	plate_dat = plate_dat.replace("DATE", ts)
	plate_dat = plate_dat.replace("TITLE", title)

	wx.write_out(wx_dir+'/plots/wx.html', plate_dat, 'w')

if __name__ == "__main__":
	i2c_addr = 0x77
	press_cal = 0.7 # kPa

	node_n = os.uname().nodename

	if node_n == 'keen':
		i2c_addr = 0x76
		title = "Indoor Weather from my Apartment in Berlin"

	if node_n == 'cutie':
		title = "Outdoor Weather from my Balcony in Berlin"

	if node_n == 'infinity':
		title = "Outdoor Weather from the west side of my Apt in Berlin"

	ts = time.strftime("%Y%m%d%H%M%S", time.gmtime())
	pi_temp = wx.pi_temp_read()
	pi_dat_string = "%s\t%s" % (ts, pi_temp)
	wx.write_out_dat_stamp(ts, 'pi_temp', pi_dat_string, wx_dir)

	(e_temp, e_hum, press, gas_r) = wx.bme680_read(i2c_addr)

	bme_dat_string = \
	"%s\tTemp: %.2f C\tHumidity: %.2f %%\tPressure: %.3f kPa\tAirQ: %d Ohms\n" \
	% (ts, e_temp, e_hum, press, gas_r)

	wx.write_out_dat_stamp(ts, 'bme680.dat', bme_dat_string, wx_dir)

	abs_hum = wx.abs_hum_g_mmm(e_temp, e_hum)
	Tdew = wx.dew_point_c(e_temp, e_hum)
	Tw = wx.web_bulb_temp(e_temp, e_hum)
	heat_i = wx.heat_index(e_temp, e_hum)
	if (heat_i == -1):
		heat_i = e_temp  # dirty hack cuz the model is pretty wonky outside a narrow range

	derived_dat_string = \
	"%s\tAbsolute Humidity: %.2f g/m³\tHeat Index: %.2f °C\tDew Point: %.2f °C\tWet Bulb Temp: %.2f °C\n" \
	% (ts, abs_hum, heat_i, Tdew, Tw)
	wx.write_out_dat_stamp(ts, 'derived.dat', derived_dat_string, wx_dir)

	gen_index(e_temp, e_hum, press + press_cal, float(pi_temp) / 1000, Tdew, abs_hum, heat_i, Tw, title)
