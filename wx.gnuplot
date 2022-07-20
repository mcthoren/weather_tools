set xtics 7200 rotate by 30 offset -5.7, -2.2
set y2tics 
set mytics
set key outside below
set xlabel "Time (UTC)" offset 0.0, -1.6;
set xdata time;
set format x "%F\n%TZ"
set timefmt "%Y%m%d%H%M%S"
set grid
set term pngcairo size 2000, 512 font ",10"

set title "Atmospheric pressure for the Last \\~48 Hours"
set ylabel "Pressure (kPa)"
set y2label "Pressure (kPa)"
set output 'pressure.png'
set format y "%.2f"
set format y2 "%.2f"

TD='/home/ghz/wx/data/bme680.dat.2-3_day'
TDPT='/home/ghz/wx/data/pi_temp.2-3_day'
TDDD='/home/ghz/wx/data/derived.dat.2-3_day'

# calibration from https://www.dwd.de/DE/wetter/wetterundklima_vorort/bayern/augsburg/_node.html [kPa]
plot TD using 1:(($9 + 5.900)) title 'Atmospheric Pressure (kPa)' with lines linecolor rgb "#A020F0"

set format y "%.1f"
set format y2 "%.1f"

set title "Raspberry Pi Temperature for the Last \\~48 Hours"
set ylabel "Pi Temp (°C)"
set y2label "Pi Temp (°C)"
set output 'pitemp.png'
plot TDPT using 1:(($2/1000)) title 'Pi Temp (°C)' with lines linecolor rgb "#00DD00" smooth bezier

set title "Temperature for the Last \\~48 Hours"
set ylabel "Temp (°C)"
set y2label "Temp (°C)"
set output 'exttemp.png'
plot TD using 1:3 title 'Temp (°C)' with lines linecolor rgb "#00DD00"

set title "Temperature, Dew Point, and Wet Bulb Temperature for the Last \\~48 Hours"
set output 'extdewtemp.png'
plot TD using 1:3 title 'Temp (°C)' with lines linecolor rgb "#00DD00", \
TDDD using 1:12 title 'Dew Point (°C)' with lines linecolor rgb "#00FFFF", \
TDDD using 1:17 title 'Wet Bulb Temp (°C)' with lines linecolor rgb "#00aaFF"

set title "Dew Point for the Last \\~48 Hours"
set ylabel "Dew Point (°C)"
set y2label "Dew Point (°C)"
set output 'extdew.png'
plot TDDD using 1:12 title 'Dew Point (°C)' with lines linecolor rgb "#00FFFF"

set title "Wet Bulb Temp for the Last \\~48 Hours"
set ylabel "Wet Bulb Temp (°C)"
set y2label "Wet Bulb Temp (°C)"
set output 'extwetbulb.png'
plot TDDD using 1:17 title 'Wet Bulb Temp (°C)' with lines linecolor rgb "#00aaff"

set title "Heat Index for the Last \\~48 Hours"
set ylabel "Heat Index (°C)"
set y2label "Heat Index (°C)"
set output 'heatindex.png'
plot TDDD using 1:8 title 'Heat Index (°C)' with lines linecolor rgb "#FF0000"

set title "Absolute Humidity for the Last \\~48 Hours"
set ylabel "Absolute Humidity (g/m³)"
set y2label "Absolute Humidity (g/m³)"
set output 'abshum.png'
plot TDDD using 1:4 title 'Absolute Humidity (g/m³)' with lines linecolor rgb "#0000DD"

set format y "%.2f"
set format y2 "%.2f"

set title "Relative Humidity for the Last \\~48 Hours"
set ylabel "Relative Humidity (%)"
set y2label "Relative Humidity (%)"
set output 'exthum.png'
plot TD using 1:6 title 'Relative Humidity (%)' with lines linecolor rgb "#00DDDD"

set format y
set format y2

set title "Air Quality Sensor Data for the Last \\~48 Hours"
set ylabel "Air Quality (kOhms)"
set y2label "Air Quality (kOhms)"
set output 'airquality.png'
plot TD using 1:(($12 / 1000)) title 'Air Quality Sensor Reading (kOhms)' with lines linecolor rgb "#FF0000"
