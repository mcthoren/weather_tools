#!/bin/bash
# see http://slackology.net/wxnotes.html for more info
# note that this is a linux version of a script normally run on BSD systems, so things like the date are a bit different.
# this is part of an ongoing project to tinker together a weather station from a raspberry pi.

OVER_DIR="/home/ghz/wx"
BASE_DIR="$OVER_DIR/plots"
LOCK=$BASE_DIR/LOCK.wx
# AVG_FILTER=$BASE_DIR/avg_filter
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

# i get so much junk in the file system every reboot with linux, that we resort to this. that data is _perfect_, or skipped.
# we start the pattern matching with 2, in the hopes that i have sth better figured out before the year 2999.
# the history of computer science suggets i won't.

# python script outputs %.2f and tabs.
PAT="^$2([0-9]{14})\t-?[0-9]*\.[0-9]{2} C\t\t[0-9]*\.[0-9]{3} kPa$"

# paste outputs a tab, the /proc file _seems_ to stick with 5 digits...
PAT1="^$2([0-9]{14})\t-?[0-9]{5}$"

# and again for the sht11. holy chrome.
# the sht11 will display humidity over 100%, since it does it when there's a lot of fog and mist around, i'm gona run with it.
PAT2="^2([0-9]{13})\tTemp: -?[0-9]?[0-9]\.[0-9]{2} C\tHumidity: [0-9]?[0-9]?[0-9]\.[0-9]{2} %\tDew Point: -?[0-9]?[0-9]\.[0-9]{2} C$"
# PAT3="^2([0-9]{13})\tTemp: -?[0-9]?[0-9]\.[0-9]{2} C\tHumidity: [0-9]?[0-9]?[0-9]\.[0-9]{2} %$"

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

PRESSURE="\"$TD.pressure\" using 1:2 title 'Pressure (kPa)' with lines linecolor rgb \"#A020F0\""
INT_TEMP="\"$TD.int_temp\" using 1:2 title 'Int Temp 0 (C)' with lines lt 1 smooth bezier"
INT_TEMP_SHT="\"$TDSHT\" using 1:3 title 'Int Temp 1 (C)' with lines linecolor rgb \"#FF00FF\""
INT_HUM_SHT="\"$TDSHT\" using 1:6 title 'Int Humidity (%)' with lines linecolor rgb \"#0000C0\" smooth bezier"
INT_DP_SHT="\"$TDSHT\" using 1:10 title 'Int Dew Point (C)' with lines linecolor rgb \"#0000FF\""
EXT_TEMP="\"$TD.ext_temp\" using 1:2 title 'Ext Temp (C)' with lines linecolor rgb \"#00DD00\""
EXT_HUM="\"$TD.ext_hum\" using 1:2 title 'Ext Humidity (%)' with lines linecolor rgb \"#00DDDD\" smooth bezier"
EXT_DP="\"$TD.ext_dp\" using 1:2 title 'Ext Dew Point (C)' with lines linecolor rgb \"#000FFF\""
PI_TEMP="\"$TDPT\" using 1:((\$2/1000)) title 'Pi Temp (C)' with lines lt 2 smooth bezier"

PROLOUGE="set title \"Weather Data for the Last ~48 Hours\";\
	set xtics 7200 ;\
	set y2tics ;\
	set key outside below;\
	set xlabel \"Time (UTC)\" offset 0.0, -1.0;\
	set xdata time;\
	set format x \"%m/%d\\n%H:%M\";\
	set timefmt \"%Y%m%d%H%M%S\";\
	set grid;\
	set term png size 2000, 512 font \",10\";"

echo "$PROLOUGE\
	set ylabel \"kPa\";\
	set y2label \"kPa\";\
	set output 'pressure.png';\
	set format y \"%.2f\" ;\
	set format y2 \"%.2f\" ;\
	plot $PRESSURE;\
	" |gnuplot

echo "$PROLOUGE\
	set ylabel \"Deg (C)\";\
	set y2label \"Deg (C)\";\
	set output 'inttemp.png';\
	plot $INT_TEMP, $INT_TEMP_SHT;\
	set output 'pitemp.png';\
	plot $PI_TEMP;\
	set output 'exttemp.png';\
	plot $EXT_TEMP, $EXT_DP;\
	" |gnuplot

#	set output 'exttemp.png';\
#	plot $EXT_TEMP;\
#	" |gnuplot

echo "$PROLOUGE\
	set ylabel \"Humidity (%)\";\
	set y2label \"Humidity (%)\";\
	set output 'inthumsht.png';\
	plot $INT_HUM_SHT;\
	set output 'exthum.png';\
	plot $EXT_HUM;\
	" |gnuplot

	# this is kinda ghetto, but to avoid race conditions we're just going to chain this together serially for now
	/home/ghz/wx/particle/graph

	# pipe 2 to /dev/null to try to quite noise caused by my horrible internet connection for now.
	# omg my intenet sux. add timeout
	rsync -e "ssh -q" --timeout=60 -ur $OVER_DIR/* wx@darkdata.org:/wx/test/ 2>/dev/null

rm $LOCK
