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

# for the lnd 712 we think the conversion factor between collected cpm data and nano grays per hours is 1000/108 for Co60.
# http://www.lndinc.com/products/348/  <--"Gamma sensitivity Co60 (cps/mr/hr)   18"
# ^--or at least that's the sense we have from the above
# which is likely going to yield a much higher dose rate than what we're actually seeing, as Co60 is prly
# not the typical isotope that we measure here in my apt.

set title "Radioactivity over the last ~8 hours."
set xtics 3600
set output 'particle_cpm_8.png'
plot OF.".8" using 1:(($2*1000/108)) title 'Dose Rate (nGy/h)' with lines lt 2,\
OF.".8.avg" using 1:(($2*1000/108)) title '16 pt Running Average Dose Rate (nGy/h)' with lines linecolor rgb "#0000d0",\
OF.".8.avg.64" using 1:(($2*1000/108)) title '64 pt Running Average Dose Rate (nGy/h)' with lines lt 1

set output 'particle_cpm_8_avg.png'
plot OF.".8.avg" using 1:(($2*1000/108)) title '16 pt Running Average Dose Rate (nGy/h)' with lines linecolor rgb "#0000d0",\
OF.".8.avg.64" using 1:(($2*1000/108)) title '64 pt Running Average Dose Rate (nGy/h)' with lines lt 1

set title "Radioactivity over the last ~24 hours."
set xtics 3600
set output 'particle_cpm_24.png'
plot OF.".24" using 1:(($2*1000/108)) title 'Dose Rate (nGy/h)' with lines lt 2,\
OF.".24.avg" using 1:(($2*1000/108)) title '16 pt Running Average Dose Rate (nGy/h)' with lines linecolor rgb "#0000d0",\
OF.".24" using 1:(($2*1000/108)) title 'Bezier Smoothed Dose Rate (nGy/h)' with lines lt 1 smooth bezier

set output 'particle_cpm_24_avg.png'
plot OF.".24.avg" using 1:(($2*1000/108)) title '16 pt Running Average Dose Rate (nGy/h)' with lines linecolor rgb "#0000d0",\
OF.".24.avg.64" using 1:(($2*1000/108)) title '64 pt Running Average Dose Rate (nGy/h)' with lines lt 1

set title "Radioactivity: Daily averages."
set xlabel "Date (yyyy/mm, UTC)" offset 0.0, -1.0
set format x "%Y/%m" 
set output 'particle_cpm_DA.png'
set xtics auto
set mxtics 2
set grid mxtics
plot DAF using 1:(($2*1000/108)) title 'Daily Average Dose Rates (nGy/h)' with histeps linecolor rgb "#0088FF"

set output 'particle_cpm_MA.png'
set timefmt "%Y%m"
set title "Radioactivity: Monthly averages."
plot MAF using 1:(($2*1000/108)) title 'Monthly Average Dose Rates (nGy/h)' with boxes lt 1

set title "Radioactivity: Daily averages."
set timefmt "%Y%m%d%H%M%S"
set grid nomxtics
set xtics 172800
set mxtics 2
set title "Radioactivity: Daily averages for the last 45 days."
set xlabel "Date (mm/dd, UTC)" offset 0.0, -1.0
set format x "%m/%d"
set output 'particle_cpm_DA.45.png'
plot DAF.".45" using 1:(($2*1000/108)) title 'Daily Average Dose Rates (nGy/h)' with boxes lt 1
