package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-co-op/gocron"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func main() {
	// Get env vars
	url, ok := os.LookupEnv("INFLUX_URL")
	if !ok {
		log.Fatalln("INFLUX_URL not set")
	}
	token, ok := os.LookupEnv("INFLUX_TOKEN")
	if !ok {
		log.Fatalln("INFLUX_TOKEN not set")
	}
	bucket, ok := os.LookupEnv("INFLUX_BUCKET")
	if !ok {
		log.Fatalln("INFLUX_BUCKET not set")
	}
	org, ok := os.LookupEnv("INFLUX_ORG")
	if !ok {
		log.Fatalln("INFLUX_ORG not set")
	}
	server, ok := os.LookupEnv("INFLUX_TAG_SERVER")
	if !ok {
		log.Fatalln("INFLUX_TAG_SERVER not set")
	}

	client := influxdb2.NewClient(url, token)
	defer client.Close()

	// run cron every 5 seconds
	s := gocron.NewScheduler(time.UTC)
	fmt.Println("Starting scheduler...")
	s.Every(5).Seconds().Do(func() {
		PublishCPU(client, bucket, org, server)
	})

	s.StartBlocking()
}
