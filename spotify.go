package main

import (
	"log"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type SpotifyStatus struct {
	ID      string
	Playing bool
}

func (s SpotifyStatus) Equal(other SpotifyStatus) bool {
	return s.ID == other.ID && s.Playing == other.Playing
}

var spotify struct {
	currentTrackSubscribers   []chan SpotifyStatus
	oldStatus                 SpotifyStatus
	currentTrackSubscribersMu sync.Mutex
	watcherOnce               sync.Once
}

func spotifyWatcher() {
	go func() {
		ticker := time.NewTicker(time.Second)
		for _ = range ticker.C {
			out, err := exec.Command("osascript", "-e", `tell application "Spotify" to player state & id of current track & name of current track & artist of current track`).Output()
			if err != nil {
				log.Printf("getting spotify current track error: %v\n", err)
				continue
			}
			fields := strings.Split(string(out), ", ")
			if len(fields) >= 2 {
				current := SpotifyStatus{Playing: fields[0] == "playing", ID: strings.TrimSpace(fields[1])}
				if !current.Equal(spotify.oldStatus) {
					if current.Playing {
						log.Printf("spotify current playing: %s by %s\n", strings.TrimSpace(fields[2]), strings.TrimSpace(fields[3]))
					} else {
						log.Printf("spotify paused/stopped\n")
					}
					(&spotify.currentTrackSubscribersMu).Lock()
					for _, ch := range spotify.currentTrackSubscribers {
						ch <- current
					}
					spotify.oldStatus = current
					(&spotify.currentTrackSubscribersMu).Unlock()
				}
			}
		}
	}()
}

func WatchSpotifyCurrentTrack() (trackIDs <-chan SpotifyStatus) {
	ch := make(chan SpotifyStatus, 8)
	(&spotify.currentTrackSubscribersMu).Lock()
	defer (&spotify.currentTrackSubscribersMu).Unlock()
	spotify.currentTrackSubscribers = append(spotify.currentTrackSubscribers, ch)
	spotify.watcherOnce.Do(spotifyWatcher)
	if spotify.oldStatus.ID != "" {
		ch <- spotify.oldStatus
	}
	return ch
}
