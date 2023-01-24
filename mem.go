package main

import (
	"context"
	"log"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/shirou/gopsutil/v3/mem"
)

func PublishMem(client influxdb2.Client, bucket, org, server string) {
	time := time.Now()
	ctx := context.Background()
	ctx2 := context.WithValue(ctx, "time", time)
	publishMemUsage(ctx2, client, bucket, org, server)
}

// publishMemUsage reads the mem usage from the system and writes it to influxdb
func publishMemUsage(ctx context.Context, client influxdb2.Client, bucket, org, server string) {
	writeAPI := client.WriteAPIBlocking(org, bucket)
	free, usage := getMemUsage()
	time := ctx.Value("time").(time.Time)
	p := influxdb2.NewPointWithMeasurement("mem").
		AddTag("server", server).
		AddField("free", free).
		AddField("usage", usage).
		SetTime(time)

	err := writeAPI.WritePoint(ctx, p)
	if err != nil {
		log.Fatalln(err)
	}
}

func getMemUsage() (uint64, float64) {
	virMem, err := mem.VirtualMemory()
	if err != nil {
		log.Fatalln(err)
	}

	return virMem.Free, virMem.UsedPercent
}
