package main

import (
	"context"
	"fmt"
	"log"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func main() {
	client := influxdb2.NewClient("influxdata-uri", "token")
	defer client.Close()

	writeAPI := client.WriteAPI("arpandas", "bucket")

	p := influxdb2.NewPointWithMeasurement("temperature").
		AddTag("location", "office").
		AddField("value", 23.5).
		SetTime(time.Now())
	writeAPI.WritePoint(p)

	writeAPI.Flush()

	queryAPI := client.QueryAPI("arpandas")
	query := `from(bucket:"bucket")|> range(start: -1h)`
	result, err := queryAPI.Query(context.Background(), query)
	if err == nil {
		for result.Next() {
			if result.TableChanged() {
				fmt.Printf("table: %s\n", result.TableMetadata().String())
			}
			fmt.Printf("value: %v\n", result.Record().Value())
		}
		if result.Err() != nil {
			log.Fatalf("query parsing error: %s\n", result.Err().Error())
		}
	} else {
		log.Fatalf("query error: %s\n", err.Error())
	}
}
