package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var echonestAPIKey string
var defaultCursorInterval float64
var silent bool

func init() {
	flag.Float64Var(&defaultCursorInterval, "i", 1, "[Optional] Default cursor interval in seconds")
	flag.StringVar(&echonestAPIKey, "k", "", "[Required] EchoNest API Key")
	flag.BoolVar(&silent, "s", false, "[Optional] Silent output")
}

func setDefaultCursorRate() {
	UpdateCursorRate(time.Duration(float64(time.Second) * defaultCursorInterval))
}

func main() {
	flag.Parse()
	if echonestAPIKey == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if silent {
		log.SetOutput(ioutil.Discard)
	}
	watcher := WatchSpotifyCurrentTrack()
	for track := range watcher {
		if track.Playing {
			bpm, err := GetSpotifyBPM(track.ID)
			if err == nil {
				interval := time.Nanosecond * time.Duration(float64(time.Minute)/bpm)
				if err = UpdateCursorRate(interval); err != nil {
					log.Printf("setting cursor rate error: %v\n", err)
				}
			} else {
				log.Printf("getting bpm error: %v\n", err)
				setDefaultCursorRate()
			}
		} else {
			setDefaultCursorRate()
		}
	}
}
