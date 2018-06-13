package main

import (
	"encoding/xml"
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

func merge() {
	d1, _ := ioutil.ReadFile("05_26_18.tcx")
	var data1 Data
	xml.Unmarshal(d1, &data1)

	d2, _ := ioutil.ReadFile("05_27_18.tcx")
	var data2 Data
	xml.Unmarshal(d2, &data2)

	// add TotalTimeSeconds, DistanceMeters
	// max of maximumSpeed
	// sum of Calories

	l1 := len(data1.Activities.Activity[0].Lap.Track.Trackpoint)
	l2 := len(data2.Activities.Activity[0].Lap.Track.Trackpoint)
	points := make([]Trackpoint, l1+l2)
	k := 0
	for _, p := range data1.Activities.Activity[0].Lap.Track.Trackpoint {
		points[k] = p
		k++
	}
	for _, p := range data2.Activities.Activity[0].Lap.Track.Trackpoint {
		points[k] = p
		k++
	}

	t1 := data1.Activities.Activity[0].Lap.TotalTimeSeconds
	t2 := data2.Activities.Activity[0].Lap.TotalTimeSeconds
	totalTime := t1 + t2

	dist1 := data1.Activities.Activity[0].Lap.DistanceMeters
	dist2 := data2.Activities.Activity[0].Lap.DistanceMeters
	distance := dist1 + dist2

	speed1 := data1.Activities.Activity[0].Lap.MaximumSpeed
	speed2 := data2.Activities.Activity[0].Lap.MaximumSpeed
	maxSpeed := math.Max(speed1, speed2)

	cal1 := data1.Activities.Activity[0].Lap.Calories
	cal2 := data2.Activities.Activity[0].Lap.Calories
	cal := cal1 + cal2

	var r Data
	r.Activities.Activity = data1.Activities.Activity
	r.Activities.Activity[0].Lap.TotalTimeSeconds = totalTime
	r.Activities.Activity[0].Lap.DistanceMeters = distance
	r.Activities.Activity[0].Lap.MaximumSpeed = maxSpeed
	r.Activities.Activity[0].Lap.Calories = cal
	r.Activities.Activity[0].Lap.Track.Trackpoint = points

	x, _ := xml.MarshalIndent(r, "", "  ")
	os.Stdout.Write(x)
}

func main() {
	merge()
}
