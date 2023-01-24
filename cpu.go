package main

import (
	"context"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/shirou/gopsutil/v3/load"
)

func PublishCPU(client influxdb2.Client, bucket, org, server string) {
	time := time.Now()
	ctx := context.Background()
	ctx2 := context.WithValue(ctx, "time", time)
	publishCpuTemp(ctx2, client, bucket, org, server)
	publishCpuLoad(ctx2, client, bucket, org, server)
}

// publishCpuTemp reads the cpu temp from the system and writes it to influxdb
func publishCpuTemp(ctx context.Context, client influxdb2.Client, bucket, org, server string) {
	writeAPI := client.WriteAPIBlocking(org, bucket)
	cpuTemp := getCpuTemp()
	time := ctx.Value("time").(time.Time)
	p := influxdb2.NewPointWithMeasurement("cpu").
		AddTag("server", server).
		AddField("temp", cpuTemp).
		SetTime(time)

	err := writeAPI.WritePoint(ctx, p)
	if err != nil {
		log.Fatalln(err)
	}
}

// publishCpuLoad reads the cpu load from the system and writes it to influxdb
func publishCpuLoad(ctx context.Context, client influxdb2.Client, bucket, org, server string) {
	writeAPI := client.WriteAPIBlocking(org, bucket)
	cpuLoad := getCpuLoad()
	time := ctx.Value("time").(time.Time)
	p := influxdb2.NewPointWithMeasurement("cpu").
		AddTag("server", server).
		AddField("load", cpuLoad).
		SetTime(time)

	err := writeAPI.WritePoint(ctx, p)
	if err != nil {
		log.Fatalln(err)
	}
}

// getCpuTemp reads the cpu temp from the system
// only to be called by the publicCpuTemp function
func getCpuTemp() int {
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

// getCpuLoad reads the cpu load from the system
// only to be called by the publicCpuLoad function
func getCpuLoad() float64 {
	avgStat, err := load.Avg()
	if err != nil {
		log.Fatalln(err)
	}
	return avgStat.Load5
}
