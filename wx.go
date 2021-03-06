#!/bin/bash
# see http://slackology.net/wxnotes.html and
# https://github.com/mcthoren/weather_tools for more info
# this is about the 4th generation weather station at this point
# this script is run from cron

OVER_DIR="/home/ghz/wx"
LOCK="$OVER_DIR/plots/LOCK.wx"
BASE_DIR="$OVER_DIR/plots"
NSAMP=4000

[ -e $LOCK ] && {
	echo "$0: lock exists" | logger
	exit 1
}

# lock is also checked for and deleted on boot, in case of a crash
touch $LOCK

$OVER_DIR/sensor_data_grab.py

cd $BASE_DIR/ ||exit 1
TD="../data/last48h"
TDPT="../data/last48h.pitemp"
TDDD="../data/last48h.derived"

DY=`date +%Y`
YDY=`date -d "-1 day" +%Y`
YYDY=`date -d "-2 day" +%Y`
DATE=`date +%Y%m%d`
YDATE=`date -d "-1 day" +%Y%m%d`
YYDATE=`date -d "-2 day" +%Y%m%d`
YYDATEH=`date -d "-2 day" +%Y%m%d%H`

# make sure our data is sth like what we expect, this has evolved from necessity.
PAT0="^2([0-9]{13})\tTemp:\ -?[0-9]?[0-9]\.[0-9]{2} C\tHumidity:\ [0-9]?[0-9]?[0-9]\.[0-9]{2} %\tPressure:\ [0-9]*\.[0-9]{3} kPa\tAirQ:\ .*Ohms$"

# 4 digits for when it gets more properly cold, 5 digits otherwise. 6 digits have not been seen.
PAT1="^$2([0-9]{14})\t-?[0-9]{4,5}$"

TDS0="cat ../data/$YYDY/bme680.dat.$YYDATE ../data/$YDY/bme680.dat.$YDATE ../data/$DY/bme680.dat.$DATE | grep -aP \"$PAT0\""
TDS1="cat ../data/$YYDY/pi_temp.$YYDATE ../data/$YDY/pi_temp.$YDATE ../data/$DY/pi_temp.$DATE | grep -aP \"$PAT1\""
TDS2="cat ../data/$YYDY/derived.dat.$YYDATE ../data/$YDY/derived.dat.$YDATE ../data/$DY/derived.dat.$DATE"

TD_DUMP0="eval $TDS0"
TD_DUMP1="eval $TDS1"
TD_DUMP2="eval $TDS2"

$TD_DUMP0 | grep -A $NSAMP $YYDATEH > $TD || $TD_DUMP0 > $TD
$TD_DUMP1 | grep -A $NSAMP $YYDATEH > $TDPT || $TD_DUMP1 > $TDPT
$TD_DUMP2 | grep -A $NSAMP $YYDATEH > $TDDD || $TD_DUMP2 > $TDDD

gnuplot -e "TD='$TD';TDPT='$TDPT';TDDD='$TDDD'" $OVER_DIR/weather_specs.gnuplot

# try not to loose 20 min of data from hard reboots
sync

rsync -e "ssh -q" --timeout=60 -ur $OVER_DIR/* wx1@dunkledaten.de:/wx1/ 2>/dev/null

rm $LOCK
