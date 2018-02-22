set title "Weather Data for the Last \\~48 Hours"
set xtics 7200
set y2tics 
set key outside below
set xlabel "Time (UTC)" offset 0.0, -1.0;
set xdata time;
set format x "%m/%d\n%H:%M"
set timefmt "%Y%m%d%H%M%S"
set grid
set term png size 2000, 512 font ",10"

set ylabel "kPa"
set y2label "kPa"
set output 'pressure.png'
set format y "%.2f"
set format y2 "%.2f"
# calibration from http://www.met.fu-berlin.de/de/wetter/ [kPa]
plot TD using 1:(($9 + 0.37)) title 'Pressure (kPa)' with lines linecolor rgb "#A020F0"

set format y
set format y2

set ylabel "Deg (C)"
set y2label "Deg (C)"
set output 'pitemp.png'
plot TDPT using 1:(($2/1000)) title 'Pi Temp (C)' with lines lt 2 smooth bezier

set output 'exttemp.png'
plot TD using 1:3 title 'Ext Temp (C)' with lines linecolor rgb "#00DD00"

set output 'extdewtemp.png'
plot TD using 1:3 title 'Ext Temp (C)' with lines linecolor rgb "#00DD00",\
TD using 1:15 title 'Ext Dew Point (C)' with lines linecolor rgb "#00FFFF"

set ylabel "Relative Humidity (%)"
set y2label "Relative Humidity (%)"
set output 'exthum.png'
plot TD using 1:6 title 'External Relative Humidity (%)' with lines linecolor rgb "#00DDDD" smooth bezier

set ylabel "Air Quality (Ohms)"
set y2label "Air Quality (Ohms)"
set output 'airquality.png'
plot TD using 1:12 title 'Air Quality Sensor Reading (Ohms)' with lines lt 1
