package main

import (
	"context"
	"log"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/shirou/gopsutil/v3/disk"
)

func PublishDisk(client influxdb2.Client, bucket, org, server string) {
	time := time.Now()
	ctx := context.Background()
	ctx2 := context.WithValue(ctx, "time", time)
	publishDiskUsage(ctx2, client, bucket, org, server)
}

func publishDiskUsage(ctx context.Context, client influxdb2.Client, bucket, org, server string) {
	writeAPI := client.WriteAPIBlocking(org, bucket)
	free, usage := getDiskUsage()
	time := ctx.Value("time").(time.Time)
	p := influxdb2.NewPointWithMeasurement("disk").
		AddTag("server", server).
		AddField("free", free).
		AddField("usage", usage).
		SetTime(time)

	err := writeAPI.WritePoint(ctx, p)
	if err != nil {
		log.Fatalln(err)
	}
}

func getDiskUsage() (uint64, float64) {
	diskUsage, err := disk.Usage("/")
	if err != nil {
		log.Fatalln(err)
	}

	return diskUsage.Free, diskUsage.UsedPercent
}
