# goMovingAverage
Moving-Average Filter for smoothing noisy data.
The code reads csv file line by line and it filters each line, concurrently.

<a href="https://www.mathworks.com/help/matlab/ref/filter.html">Reference</a>

# Usage
go run main.go -windowSize=30 -dataPath="rawSensorData.csv"
