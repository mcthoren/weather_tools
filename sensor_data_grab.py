#!/usr/bin/python


import time
import datetime
import bme680

wx_dir = "/home/ghz/wx"

def write_out(file_name, data, mode):
	out_file_fd = open(file_name, mode)
	out_file_fd.write(data)
	out_file_fd.close()

def write_out_dat_stamp(ts, n_plate, data):
	f_ts = ts[0:8]
	write_out(wx_dir+'/data/'+n_plate+'.'+f_ts, data, 'a')

def bme680_read():
	sensor = bme680.BME680(i2c_addr=0x77)

	sensor.set_humidity_oversample(bme680.OS_2X)
	sensor.set_pressure_oversample(bme680.OS_4X)
	sensor.set_temperature_oversample(bme680.OS_8X)
	sensor.set_filter(bme680.FILTER_SIZE_3)
	sensor.set_gas_status(bme680.ENABLE_GAS_MEAS)

	sensor.set_gas_heater_temperature(320)
	sensor.set_gas_heater_duration(150)
	sensor.select_gas_heater_profile(0)

	ts = datetime.datetime.fromtimestamp(time.time()).strftime("%Y%m%d%H%M%S")

	temp = sensor.data.temperature
	pressure = sensor.data.pressure
	hum = sensor.data.humidity

	gas_res = sensor.data.gas_resistance

	dat_string = "%s\tTemp: %.2f C\tHumidity: %.2f %%\tPressure: %.2f kPa\tAirQ: %.0f Ohms\n" % (ts, temp, hum, pressure / 10, gas_res)

	write_out_dat_stamp(ts, 'bme680.dat', dat_string)
		
def pi_temp_read():
	temp_file = "/sys/class/thermal/thermal_zone0/temp"
	temp_file_fd = open(temp_file, 'r')

	ts = datetime.datetime.fromtimestamp(time.time()).strftime("%Y%m%d%H%M%S")

	temp_data = temp_file_fd.read()
	temp_file_fd.close()

	dat_string = "%s\t%s" % (ts, temp_data)

	write_out_dat_stamp(ts, 'pi_temp', dat_string)

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

	write_out(wx_dir+'/plots/wx.html', plate_dat, 'w')

if __name__ == "__main__":
	pi_temp = pi_temp_read()
	# gen_index(e_temp, e_hum, e_dp, press, bmp_temp, i_temp, i_hum, i_dp, pi_temp)
	bme680_read()
