package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

type Config struct {
	WindowSize int
	B          []float64
	A          []float64
}

var conf Config

func main() {

	// Getting terminal arguments with flag package and parsing them.
	ws := flag.Int("windowSize", 70, "-windowSize = [INTEGER]")
	df := flag.String("dataPath", "", "-dataPath = [PATH]")
	flag.Parse()

	b := make([]float64, int(*ws))
	for i := range b {
		b[i] = 1 / float64(*ws)
	}

	// Initializing the config
	conf = Config{
		WindowSize: *ws,
		B:          b,
		A:          []float64{1},
	}

	// Reading target file. If it fails, exit the program.
	csvFile, err := os.Open(*df)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Closing the opened .csv file with deferred self executing function.
	defer func() {
		if err := csvFile.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	// Declaring a reader for .cvs file.
	reader := csv.NewReader(bufio.NewReader(csvFile))

	// Every row that extracted from the .csv file will be computed concurrently.
	c := make(chan []float64)

	// Below loop goes until reaching the end of file (EOF)
	// It reads the file line by line, and handles the lines by calling
	// parseRows() function.
	i := 1
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
			return
		}

		go parseRows(line, c, i)
		i++
	}

	for k := 1; k < i; k++ {
		data := <-c
		fmt.Println(data)
	}
	close(c)

	fmt.Println("About to exit")
}

// This function converts the read string lines []string to floating-point values []float64 and
// calls the MovingAverage() function.
func parseRows(row []string, c chan<- []float64, line int) {
	fmt.Printf("Goroutine#%d doing its job\n", line)
	out := make([]float64, len(row))

	for i, v := range row {
		fv, _ := strconv.ParseFloat(v, 64)
		out[i] = fv
	}

	conf.MovingAverage(&out)

	fmt.Printf("Goroutine #%d has finished its job\n", line)
	c <- out
}

func (c Config) MovingAverage(d *[]float64) {

	for i, v := range *d {

		if i == 0 {
			(*d)[0] = c.B[0] * v
		} else {
			if i < len(c.B) {
				(*d)[i] += c.B[i] * v
			}
			for j := 0; j < i; j++ {
				k := i - j
				if (k < len(c.B)) && (j < len(*d)) {
					(*d)[i] += c.B[k] * v
				}
				if (k < len(*d)) && (j < len(c.A)) {
					(*d)[i] -= c.A[j] * (*d)[k]
				}
			}

		}
	}

}
