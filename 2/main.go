package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/prometheus/prometheus/prompb"
)

var queue Queue

func main() {
	metricUrl := os.Getenv("METRIC_URL")
	prometeusUrl := os.Getenv("PROMETHEUS_URL")

	getMetricTicker := time.NewTicker(2 * time.Second)
	getMetricTickerDone := make(chan struct{})

	go func() {
		for {
			select {
			case <-getMetricTicker.C:
				getMetricTask(metricUrl)
			case <-getMetricTickerDone:
				getMetricTicker.Stop()
				return
			}
		}
	}()

	sendMetricTicker := time.NewTicker(2 * time.Second)
	sendMetricTickerDone := make(chan struct{})

	go func() {
		for {
			select {
			case <-sendMetricTicker.C:
				sendMetricTask(prometeusUrl)
			case <-sendMetricTickerDone:
				sendMetricTicker.Stop()
				return
			}
		}
	}()

	time.Sleep(100 * time.Second)
}

func getMetricTask(metricUrl string) {
	metrics := getMetrics(metricUrl)
	req := serializeMetric(metrics)

	data, err := proto.Marshal(req)

	if err != nil {
		fmt.Println(err)
	}

	queue.Push(data)
}

func serializeMetric(map[string]float32) *prompb.WriteRequest {
	// did not realized
	return &prompb.WriteRequest{}
}

func sendMetricTask(prometeusUrl string) {
	data := queue.Pull()
	sendMetrics(prometeusUrl, &data)
}

func getMetrics(metricUrl string) map[string]float32 {
	resp, err := http.Get(metricUrl)

	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	responseString := string(responseData)
	metrics := strings.Split(responseString, "\n")

	result := make(map[string]float32)

	for _, metric := range metrics {
		tokens := strings.Split(metric, " ")

		value, err := strconv.ParseFloat(tokens[1], 32)

		if err != nil {
			fmt.Println(err)
		}

		result[tokens[0]] = float32(value)
	}

	return result
}

func sendMetrics(prometeusUrl string, data *[]byte) {
	_, err := http.Post(prometeusUrl, "application/x-protobuf", bytes.NewReader(*data))

	if err != nil {
		fmt.Println(err)
	}
}
