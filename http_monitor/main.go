package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

const (
	dbName = "test"
)

// Type for health check lul
type Health struct {
	status   int
	url      string
	duration float64
}

func main() {
	//
	urls := [...]string{"http://google.com",
		"http://yahoo.com",
		"http://blub.sa"}
	// HTTP client
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://localhost:8086",
	})
	if err != nil {
		log.Fatalln(err)
	}
	// point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  dbName,
		Precision: "s",
	})
	if err != nil {
		log.Fatalln(err)
	}
	tags := map[string]string{"http": "status"}
	fields := map[string]interface{}{
		"statuscode": 200,
		"rtt":        22.222,
	}
	p, err := client.NewPoint("http-check", tags, fields, time.Now())
	if err != nil {
		log.Fatalln(err)
	}
	bp.AddPoint(p)
	if err := c.Write(bp); err != nil {
		log.Fatalln(err)
	}

	// channels
	ch := make(chan Health)
	for _, url := range urls {
		go MakeRequest(url, ch)
	}

	c.Ping(1)

	for _, _ = range urls {
		fmt.Println(<-ch)
	}
}

func MakeRequest(url string, ch chan<- Health) {
	start := time.Now()
	resp, err := http.Head(url)
	if err != nil {
		sec := time.Since(start).Seconds()
		ch <- Health{
			url:      url,
			status:   404,
			duration: sec,
		}
		return
	}
	sec := time.Since(start).Seconds() / time.Millisecond.Seconds()
	ch <- Health{
		url:      url,
		status:   resp.StatusCode,
		duration: sec,
	}
}
