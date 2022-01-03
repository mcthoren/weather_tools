This code can be found in the following places:
* https://wx1.slackology.net/plots/wx.html	<--page
* https://github.com/mcthoren/weather_tools	<--code
* https://wx1.slackology.net/			<--code, data, plots, page

Much of the data for all this is currently collected from a BME680, the datasheet for which can be found here:
* https://ae-bst.resource.bosch.com/media/_tech/media/datasheets/BST-BME680-DS001.pdf

Notes:

* This is part of an ongoing project to tinker together a weather station from a raspberry pi.

* In theory, if one wanted to run this themselves, they could change the WT_DIR variable, and some of the upload paths, and be off to a good start

* The i2c bus needs to be enabled.  This is probably most easily done with raspi-config.  Make sure your user is in the i2c and gpio groups 

* apt install python3-smbus python3-bme680 #  <-- something like that will help one get off to a good start

* The sensor data used to be pretty noisy, this combined with some of the things that happen to the linux file system on the raspberry pi necessitated the very strict pattern matching observed toward the top of wx.go.

* The Air Quality readings are currently raw data from the bme680, the datasheet for which is linked above.
