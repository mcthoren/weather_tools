#!/bin/bash
# see http://slackology.net/wxnotes.html and
# https://github.com/mcthoren/weather_tools for more info
# this is about the 4th generation weather station at this point

OVER_DIR="/home/ghz/wx"
BASE_DIR="$OVER_DIR/plots"
SPIKE_FILTER=$OVER_DIR/spike_filter
NSAMP=4000

cd $BASE_DIR/ ||exit 1
TD="../data/last48h"
TDPT="../data/last48h.pitemp"

DY=`date +%Y`
DATE=`date +%Y%m%d`
YDATE=`date -d "-1 day" +%Y%m%d`
YYDATE=`date -d "-2 day" +%Y%m%d`
YYDATEH=`date -d "-2 day" +%Y%m%d%H`

# make sure our data is sth like what we expect, this has evolved from necessity.
PAT0="^2([0-9]{13})\tTemp:\ -?[0-9]?[0-9]\.[0-9]{2} C\tHumidity:\ [0-9]?[0-9]?[0-9]\.[0-9]{2} %\tPressure:\ [0-9]*\.[0-9]{3} kPa\tAirQ:\ .*Ohms$"

# paste outputs a tab, the /proc file _seems_ to stick with 5 digits...
PAT1="^$2([0-9]{14})\t-?[0-9]{5}$"

TDS0="cat ../data/bme680.dat.$YYDATE ../data/bme680.dat.$YDATE ../data/bme680.dat.$DATE | grep -aP \"$PAT0\""
TDS1="cat ../data/pi_temp.$YYDATE ../data/pi_temp.$YDATE ../data/pi_temp.$DATE | grep -aP \"$PAT1\""

TD_DUMP0="eval $TDS0"
TD_DUMP1="eval $TDS1"

$TD_DUMP0 | grep -A $NSAMP $YYDATEH > $TD || $TD_DUMP0 > $TD
$TD_DUMP1 | grep -A $NSAMP $YYDATEH > $TDPT || $TD_DUMP1 > $TDPT

# was broken out into individual files to assist possible future filtering (bogus values, averaging, etc)
# we're going to try to skip this part and see how it goes
# $TD_DUMP0 |awk '{print $1, $4}' | $SPIKE_FILTER | grep -A $NSAMP $YYDATEH > $TD.pressure || $TD_DUMP0 |awk '{print $1, $4}' | $SPIKE_FILTER > $TD.pressure

# cat $TD_EXT | awk '{print $1, $3}' | $SPIKE_FILTER | grep -A $NSAMP $YYDATEH > $TD.ext_temp || cat $TD_EXT | awk '{print $1, $3}' | $SPIKE_FILTER > $TD.ext_temp
# cat $TD_EXT | awk '{print $1, $6}' | $SPIKE_FILTER | grep -A $NSAMP $YYDATEH > $TD.ext_hum || cat $TD_EXT | awk '{print $1, $6}' | $SPIKE_FILTER > $TD.ext_hum
# cat $TD_EXT | awk '{print $1, $10}' | $SPIKE_FILTER | grep -A $NSAMP $YYDATEH > $TD.ext_dp || cat $TD_EXT | awk '{print $1, $10}' | $SPIKE_FILTER > $TD.ext_dp

# gnuplot -e "TD='$TD';TDPT='$TDPT'" $OVER_DIR/weather_specs.gnuplot
