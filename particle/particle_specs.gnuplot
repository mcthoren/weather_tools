set xtics
set y2tics
set key outside below
set xlabel "Time (UTC)" offset 0.0, -1.0
set xdata time
set format x "%m/%d\n%H:%M"
set timefmt "%Y%m%d%H%M%S"
set grid
set ylabel "nGy/h"
set y2label "nGy/h"
set term png size 2000, 512 font ",10"

# CPM8="\"$OF.8\" using 1:((\$2*1000/108)) title 'Dose Rate (nGy/h)' with lines lt 2"
# AVG_CPM8="\"$OF.8.avg\" using 1:((\$2*1000/108)) title '16 pt Running Average Dose Rate (nGy/h)' with lines linecolor rgb \"#0000d0\""
# AVG_648="\"$OF.8.avg.64\" using 1:((\$2*1000/108)) title '64 pt Running Average Dose Rate (nGy/h)' with points lt 1"
# AVG_648l="\"$OF.8.avg.64\" using 1:((\$2*1000/108)) title '64 pt Running Average Dose Rate (nGy/h)' with lines lt 1"

# BEZ8="\"$OF.8\" using 1:((\$2*1000/108)) title 'Bezier Smoothed Dose Rate (nGy/h)' with lines lt 1 smooth bezier"
# BEZ24="\"$OF.24\" using 1:((\$2*1000/108)) title 'Bezier Smoothed Dose Rate (nGy/h)' with lines lt 1 smooth bezier"

# CPM24="\"$OF.24\" using 1:((\$2*1000/108)) title 'Dose Rate (nGy/h)' with lines lt 2"
# AVG_CPM24="\"$OF.24.avg\" using 1:((\$2*1000/108)) title '16 pt Running Average Dose Rate (nGy/h)' with lines linecolor rgb \"#0000d0\""
# AVG_6424="\"$OF.24.avg.64\" using 1:((\$2*1000/108)) title '64 pt Running Average Dose Rate (nGy/h)' with points lt 1"
# AVG_6424l="\"$OF.24.avg.64\" using 1:((\$2*1000/108)) title '64 pt Running Average Dose Rate (nGy/h)' with lines lt 1"

# DA_BAR="\"$DA_OF\" using 1:((\$2*1000/108)) title 'Daily Average Dose Rates (nGy/h)' with histep lt 3"
# DA_BAR_45="\"$DA_OF.45\" using 1:((\$2*1000/108)) title 'Daily Average Dose Rates (nGy/h)' with boxes lt 1"
# MA_BAR="\"$MA_OF\" using 1:((\$2*1000/108)) title 'Monthly Average Dose Rates (nGy/h)' with boxes lt 1"

set title "Radioactivity over the last ~8 hours."
set xtics 3600
set output '$PLOT_DIR/particle_cpm_8.png'
# plot $CPM8, $AVG_CPM8, $AVG_648l
plot OF.".8" using 1:(($2*1000/108)) title 'Dose Rate (nGy/h)' with lines lt 2
plot OF.".8.avg" using 1:(($2*1000/108)) title '16 pt Running Average Dose Rate (nGy/h)' with lines linecolor rgb "#0000d0"
plot OF.".8.avg.64" using 1:(($2*1000/108)) title '64 pt Running Average Dose Rate (nGy/h)' with lines lt 1

set output '$PLOT_DIR/particle_cpm_8_avg.png'
#plot $AVG_CPM8, $AVG_648l
plot OF.".8.avg" using 1:(($2*1000/108)) title '16 pt Running Average Dose Rate (nGy/h)' with lines linecolor rgb "#0000d0"
plot OF.".8.avg.64" using 1:(($2*1000/108)) title '64 pt Running Average Dose Rate (nGy/h)' with lines lt 1

set title "Radioactivity over the last ~24 hours."
set xtics 3600
set output '$PLOT_DIR/particle_cpm_24.png'
# plot $CPM24, $AVG_CPM24, $BEZ24
plot OF.".24" using 1:(($2*1000/108)) title 'Dose Rate (nGy/h)' with lines lt 2
plot OF.".24.avg" using 1:(($2*1000/108)) title '16 pt Running Average Dose Rate (nGy/h)' with lines linecolor rgb "#0000d0"
plot OF.".24" using 1:(($2*1000/108)) title 'Bezier Smoothed Dose Rate (nGy/h)' with lines lt 1 smooth bezier

set output '$PLOT_DIR/particle_cpm_24_avg.png'
# plot $AVG_CPM24, $AVG_6424l
plot OF.".24.avg" using 1:(($2*1000/108)) title '16 pt Running Average Dose Rate (nGy/h)' with lines linecolor rgb "#0000d0"
plot OF.".24.avg.64" using 1:(($2*1000/108)) title '64 pt Running Average Dose Rate (nGy/h)' with lines lt 1

set title "Radioactivity: Daily averages."
set xlabel "Date (yyyy/mm, UTC)" offset 0.0, -1.0
set format x "%Y/%m" 
set output '$PLOT_DIR/particle_cpm_DA.png'
set mxtics 2
set grid mxtics
# plot $DA_BAR
plot DA_OF.".45" using 1:(($2*1000/108)) title 'Daily Average Dose Rates (nGy/h)' with boxes lt 1

set output '$PLOT_DIR/particle_cpm_MA.png'
set timefmt "%Y%m"
set title "Radioactivity: Monthly averages."
# plot $MA_BAR
plot MA_OF using 1:(($2*1000/108)) title 'Monthly Average Dose Rates (nGy/h)' with boxes lt 1

set title "Radioactivity: Daily averages."
set timefmt "%Y%m%d%H%M%S"
set grid nomxtics
set xtics 172800
set mxtics 2
set title "Radioactivity: Daily averages for the last 45 days."
set xlabel "Date (mm/dd, UTC)" offset 0.0, -1.0
set format x "%m/%d"
set output '$PLOT_DIR/particle_cpm_DA.45.png'
# plot $DA_BAR_45
plot DA_OF.".45" using 1:(($2*1000/108)) title 'Daily Average Dose Rates (nGy/h)' with boxes lt 1
