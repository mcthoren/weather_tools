#!/bin/bash
# see http://slackology.net/wxnotes.html for more info

OVER_DIR="/home/ghz/wx"
BASE_DIR="$OVER_DIR/plots"
LOCK=$BASE_DIR/LOCK.wx
SPIKE_FILTER=$OVER_DIR/spike_filter
NSAMP=4000

cd $BASE_DIR/ ||exit 1
TD="../data/last48h"
TDPT="../data/last48h.pitemp"
TDSHT="../data/last48h.sht"
TD_EXT="../data/last48h.htu_ext"

[ -e $LOCK ] && {
	echo "$0: lock exists" | logger
	exit 1
}

# lock is also checked for and deleted on boot, in case of a crash
touch $LOCK

DY=`date +%Y`
DATE=`date +%Y%m%d`
YDATE=`date -d "-1 day" +%Y%m%d`
YYDATE=`date -d "-2 day" +%Y%m%d`
YYDATEH=`date -d "-2 day" +%Y%m%d%H`

# python script outputs %.2f and tabs.
PAT="^$2([0-9]{14})\t-?[0-9]*\.[0-9]{2} C\t\t[0-9]*\.[0-9]{3} kPa$"

# paste outputs a tab, the /proc file _seems_ to stick with 5 digits...
PAT1="^$2([0-9]{14})\t-?[0-9]{5}$"
PAT2="^2([0-9]{13})\tTemp: -?[0-9]?[0-9]\.[0-9]{2} C\tHumidity: [0-9]?[0-9]?[0-9]\.[0-9]{2} %\tDew Point: -?[0-9]?[0-9]\.[0-9]{2} C$"

TDS0="cat ../data/bmp0085_grab.dat.$YYDATE ../data/bmp0085_grab.dat.$YDATE ../data/bmp0085_grab.dat.$DATE | grep -aP \"$PAT\""
TDS1="cat ../data/pi_temp.$YYDATE ../data/pi_temp.$YDATE ../data/pi_temp.$DATE | grep -aP \"$PAT1\""
TDS2="cat ../data/sht11_grab.dat.$YYDATE ../data/sht11_grab.dat.$YDATE ../data/sht11_grab.dat.$DATE | grep -aP \"$PAT2\""
TDS3="cat ../data/htu21df_grab.dat.$YYDATE ../data/htu21df_grab.dat.$YDATE ../data/htu21df_grab.dat.$DATE | grep -aP \"$PAT2\""

TD_DUMP0="eval $TDS0"
TD_DUMP1="eval $TDS1"
TD_DUMP2="eval $TDS2"
TD_DUMP3="eval $TDS3"

$TD_DUMP0 | grep -A $NSAMP $YYDATEH > $TD || $TD_DUMP0 > $TD
$TD_DUMP1 | grep -A $NSAMP $YYDATEH > $TDPT || $TD_DUMP1 > $TDPT
$TD_DUMP2 | grep -A $NSAMP $YYDATEH > $TDSHT || $TD_DUMP2 > $TDSHT
$TD_DUMP3 | grep -A $NSAMP $YYDATEH > $TD_EXT || $TD_DUMP3 > $TD_EXT

# broken out into individual files to assist possible future filtering (bogus values, averaging, etc)
$TD_DUMP0 |awk '{print $1, $2}' | grep -A $NSAMP $YYDATEH > $TD.int_temp || $TD_DUMP0 |awk '{print $1, $2}' > $TD.int_temp
$TD_DUMP0 |awk '{print $1, $4}' | grep -A $NSAMP $YYDATEH > $TD.pressure || $TD_DUMP0 |awk '{print $1, $4}' > $TD.pressure

cat $TD_EXT | awk '{print $1, $3}' | $SPIKE_FILTER | grep -A $NSAMP $YYDATEH > $TD.ext_temp || cat $TD_EXT | awk '{print $1, $3}' | $SPIKE_FILTER > $TD.ext_temp
cat $TD_EXT | awk '{print $1, $6}' | $SPIKE_FILTER | grep -A $NSAMP $YYDATEH > $TD.ext_hum || cat $TD_EXT | awk '{print $1, $6}' | $SPIKE_FILTER > $TD.ext_hum
cat $TD_EXT | awk '{print $1, $10}' | $SPIKE_FILTER | grep -A $NSAMP $YYDATEH > $TD.ext_dp || cat $TD_EXT | awk '{print $1, $10}' | $SPIKE_FILTER > $TD.ext_dp

$BASE_DIR/../gen_index

gnuplot -e "TD='$TD';TDPT='$TDPT';TDSHT='$TDSHT'" $OVER_DIR/weather_specs.gnuplot

# this is kinda ghetto, but to avoid race conditions we're just going to chain this together serially for now
/home/ghz/wx/particle/graph

rsync -e "ssh -q" --timeout=60 -ur $OVER_DIR/* wx@darkdata.org:/wx/test/ 2>/dev/null

rm $LOCK
