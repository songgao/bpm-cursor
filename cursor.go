package main

import (
	"log"
	"os/exec"
	"strconv"
	"time"
)

func updateCursorRateITerm3(duration time.Duration) error {
	return exec.Command("defaults", "write", "com.googlecode.iterm2", "TimeBetweenBlinks", "-float", strconv.FormatFloat(duration.Seconds(), 'f', -1, 64)).Run()
}

func UpdateCursorRate(duration time.Duration) error {
	log.Printf("updating cursor interval: %s\n", duration.String())
	return updateCursorRateITerm3(duration)
}
