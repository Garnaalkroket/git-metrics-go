package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

type repoTime struct {
	name string
	time float64
}

//Define a struct for you collector that contains pointers
//to prometheus descriptors for each metric you wish to expose.
//Note you can also include fields of other types if they provide utility
//but we just won't be exposing them as metrics.
type gitCollector struct {
	gitMetric *prometheus.Desc
	gitRegExp *regexp.Regexp
}

//You must create a constructor for you collector that
//initializes every descriptor and returns a pointer to the collector
func newGitCollector(r *regexp.Regexp) *gitCollector {
	return &gitCollector{
		gitMetric: prometheus.NewDesc("git_metric",
			"Shows whether a foo has occurred in our cluster",
			[]string{"repository"}, nil,
		),
		gitRegExp: r,
	}
}

//Each and every collector must implement the Describe function.
//It essentially writes all descriptors to the prometheus desc channel.
func (collector *gitCollector) Describe(ch chan<- *prometheus.Desc) {

	//Update this section with the each metric you create for a given collector
	ch <- collector.gitMetric
}

//Collect implements required collect function for all promehteus collectors
func (collector *gitCollector) Collect(ch chan<- prometheus.Metric) {
	my_time := readLines(*collector.gitRegExp)
	//Implement logic here to determine proper metric value to return to prometheus
	//for each descriptor or call other functions that do so.

	//Write latest value for each metric in the prometheus metric channel.
	//Note that you can pass CounterValue, GaugeValue, or UntypedValue types here.
	ch <- prometheus.MustNewConstMetric(collector.gitMetric, prometheus.GaugeValue, my_time.time, my_time.name)

}

func readLines(r regexp.Regexp) repoTime {
	file_location := os.Getenv("LOG_FILE_LOCATION")
	fmt.Printf("Parsing latest entries of %s", file_location)
	file, err := os.Open(file_location)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	buf := make([]byte, 500)
	stat, err := os.Stat(file_location)
	start := stat.Size() - 500
	fmt.Printf("\nCurrent size %d, reading %d starting from %d", stat.Size(), len(buf), start)
	_, err = file.ReadAt(buf, start-1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Last entries: %s\n", buf)
	fmt.Printf("Matching against: %s", r.String())
	matches := r.FindAll(buf, -1)
	repo := "NA"
	time_float := 0.00
	if matches != nil {
		m := string(matches[len(matches)-1])
		print("Last match: ")
		print(m)
		parts := strings.Split(m, " ")
		repo = parts[4]
		time := parts[2]
		time = strings.ReplaceAll(time, "'", "")
		time_float, err = strconv.ParseFloat(time, 64)
		if err != nil {

		}
		fmt.Printf("%s %.2f", repo, time_float)
	}
	return repoTime{repo, time_float}
}
