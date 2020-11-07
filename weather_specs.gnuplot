set title "Weather Data for the Last \\~48 Hours"
set xtics 7200 rotate by 30 offset -5.7, -2.2
set y2tics 
set mytics
set key outside below
set xlabel "Time (UTC)" offset 0.0, -1.6;
set xdata time;
set format x "%F\n%TZ"
set timefmt "%Y%m%d%H%M%S"
set grid
set term png size 2000, 512 font ",10"

set ylabel "kPa"
set y2label "kPa"
set output 'pressure.png'
set format y "%.2f"
set format y2 "%.2f"
# calibration from https://www.dwd.de/DE/wetter/wetterundklima_vorort/bayern/augsburg/_node.html [kPa]
plot TD using 1:(($9 + 5.900)) title 'Pressure (kPa)' with lines linecolor rgb "#A020F0"

set format y "%.1f"
set format y2 "%.1f"

set ylabel "(°C)"
set y2label "(°C)"
set output 'pitemp.png'
plot TDPT using 1:(($2/1000)) title 'Pi Temp (°C)' with lines linecolor rgb "#00DD00" smooth bezier

set output 'exttemp.png'
plot TD using 1:3 title 'Ext Temp (°C)' with lines linecolor rgb "#00DD00"

set output 'extdewtemp.png'
plot TD using 1:3 title 'Ext Temp (°C)' with lines linecolor rgb "#00DD00",\
TDDD using 1:12 title 'Ext Dew Point (°C)' with lines linecolor rgb "#00FFFF"

set output 'extdew.png'
plot TDDD using 1:12 title 'Ext Dew Point (°C)' with lines linecolor rgb "#00FFFF"

set output 'heatindex.png'
plot TDDD using 1:8 title 'Heat Index (°C)' with lines linecolor rgb "#FF0000"

set ylabel "Absolute Humidity (g/m³)"
set y2label "Absolute Humidity (g/m³)"
set output 'abshum.png'
plot TDDD using 1:4 title 'Absolute Humidity (g/m³)' with lines linecolor rgb "#0000DD"

set format y "%.2f"
set format y2 "%.2f"

set ylabel "Relative Humidity (%)"
set y2label "Relative Humidity (%)"
set output 'exthum.png'
plot TD using 1:6 title 'External Relative Humidity (%)' with lines linecolor rgb "#00DDDD" smooth bezier

set format y
set format y2

set ylabel "Air Quality (kOhms)"
set y2label "Air Quality (kOhms)"
set output 'airquality.png'
plot TD using 1:(($12 / 1000)) title 'Air Quality Sensor Reading (kOhms)' with lines linecolor rgb "#FF0000"
