#!/bin/bash
# see http://slackology.net/wxnotes.html and
# https://github.com/mcthoren/weather_tools for more info
# this is about the 4th generation weather station at this point

OVER_DIR="/home/ghz/wx"
BASE_DIR="$OVER_DIR/plots"
NSAMP=4000

cd $BASE_DIR/ ||exit 1
TD="../data/last48h"
TDPT="../data/last48h.pitemp"

DY=`date +%Y`
YDY=`date -d "-1 day" +%Y`
YYDY=`date -d "-2 day" +%Y`
DATE=`date +%Y%m%d`
YDATE=`date -d "-1 day" +%Y%m%d`
YYDATE=`date -d "-2 day" +%Y%m%d`
YYDATEH=`date -d "-2 day" +%Y%m%d%H`

# make sure our data is sth like what we expect, this has evolved from necessity.
# PAT0="^2([0-9]{13})\tTemp:\ -?[0-9]?[0-9]\.[0-9]{2} C\tHumidity:\ [0-9]?[0-9]?[0-9]\.[0-9]{2} %\tPressure:\ [0-9]*\.[0-9]{3} kPa\tAirQ:\ .*Ohms$"

# hack the expression to match forms both with and without dewpoint for the next 48 hours
PAT0="^2([0-9]{13})\tTemp:\ -?[0-9]?[0-9]\.[0-9]{2} C\tHumidity:\ [0-9]?[0-9]?[0-9]\.[0-9]{2} %\tPressure:\ [0-9]*\.[0-9]{3} kPa\tAirQ:\ .*$"

# paste outputs a tab, the /proc file _seems_ to stick with 5 digits...
PAT1="^$2([0-9]{14})\t-?[0-9]{5}$"

TDS0="cat ../data/$DY/bme680.dat.$YYDATE ../data/$YDY/bme680.dat.$YDATE ../data/$YYDY/bme680.dat.$DATE | grep -aP \"$PAT0\""
TDS1="cat ../data/$DY/pi_temp.$YYDATE ../data/$YDY/pi_temp.$YDATE ../data/$YYDY/pi_temp.$DATE | grep -aP \"$PAT1\""

TD_DUMP0="eval $TDS0"
TD_DUMP1="eval $TDS1"

$TD_DUMP0 | grep -A $NSAMP $YYDATEH > $TD || $TD_DUMP0 > $TD
$TD_DUMP1 | grep -A $NSAMP $YYDATEH > $TDPT || $TD_DUMP1 > $TDPT

gnuplot -e "TD='$TD';TDPT='$TDPT'" $OVER_DIR/weather_specs.gnuplot
