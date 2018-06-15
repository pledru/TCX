package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math"
	"os"
)

type Trackpoint struct {
	Time     string
	Position struct {
		LatitudeDegrees  float64
		LongitudeDegrees float64
	}
	AltitudeMeters float64
	DistanceMeters float64
}

type Data struct {
	Activities struct {
		Activity []struct {
			Sport string `xml:"Sport,attr"`
			Id    string
			Lap   struct {
				StartTime        string `xml:"StartTime,attr"`
				TotalTimeSeconds int64
				DistanceMeters   float64
				MaximumSpeed     float64
				Calories         int64
				Intensity        string
				TriggerMethod    string
				Track            struct {
					Trackpoint []Trackpoint
				}
			}
		}
	}
}

func merge(files []string) {
	data := make([]Data, len(files))
	for i, f := range files {
		d, err := ioutil.ReadFile(f)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		xml.Unmarshal(d, &data[i])
	}
	l := 0
	for _, d := range data {
		l += len(d.Activities.Activity[0].Lap.Track.Trackpoint)
	}
	points := make([]Trackpoint, l)
	k := 0
	var totalTime int64
	var distance float64
	var maxSpeed float64
	var cal int64
	for _, d := range data {
		for _, p := range d.Activities.Activity[0].Lap.Track.Trackpoint {
			points[k] = p
			k++
		}
		totalTime += d.Activities.Activity[0].Lap.TotalTimeSeconds
		distance += d.Activities.Activity[0].Lap.DistanceMeters
		maxSpeed = math.Max(maxSpeed, d.Activities.Activity[0].Lap.MaximumSpeed)
		cal += d.Activities.Activity[0].Lap.Calories
	}

	var r Data
	r.Activities.Activity = data[0].Activities.Activity
	r.Activities.Activity[0].Lap.TotalTimeSeconds = totalTime
	r.Activities.Activity[0].Lap.DistanceMeters = distance
	r.Activities.Activity[0].Lap.MaximumSpeed = maxSpeed
	r.Activities.Activity[0].Lap.Calories = cal
	r.Activities.Activity[0].Lap.Track.Trackpoint = points

	x, _ := xml.MarshalIndent(r, "", "  ")
	os.Stdout.Write(x)
}

func main() {
	arg := os.Args[1:]
	merge(arg)
}
