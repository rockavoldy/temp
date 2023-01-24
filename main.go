package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
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
		publishCpuTemp(client, bucket, org, server)
	})

	s.StartBlocking()
}

func publishCpuTemp(client influxdb2.Client, bucket, org, server string) {
	// read sensor from system, and write to influxdb
	writeAPI := client.WriteAPIBlocking(org, bucket)
	cpuTemp := GetCpuTemp()
	p := influxdb2.NewPointWithMeasurement("cpu").
		AddTag("server", server).
		AddField("temp", cpuTemp).
		SetTime(time.Now())

	err := writeAPI.WritePoint(context.Background(), p)
	if err != nil {
		log.Fatalln(err)
	}
}

func GetCpuTemp() int {
	cmd := "cat /sys/class/thermal/thermal_zone0/temp"
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		log.Fatalln(err)
	}
	output := strings.TrimSuffix(string(out), "\n")
	outputInt, err := strconv.Atoi(output)
	if err != nil {
		log.Fatalln(err)
	}
	return outputInt / 1000.0
}
