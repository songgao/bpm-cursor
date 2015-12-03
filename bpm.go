package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

func GetSpotifyBPM(spotifyID string) (bpm float64, err error) {
	u, _ := url.Parse("http://developer.echonest.com/api/v4/track/profile")
	q := u.Query()
	q.Set("format", "json")
	q.Set("bucket", "audio_summary")
	q.Set("id", spotifyID)
	q.Set("api_key", echonestAPIKey)
	u.RawQuery = q.Encode()
	var resp *http.Response
	resp, err = http.Get(u.String())
	if err != nil {
		return
	}
	var data struct {
		Response struct {
			Track struct {
				AudioSummary struct {
					Tempo float64 `json:"tempo"`
				} `json:"audio_summary"`
			} `json:"track"`
			Status struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
			} `json:"status"`
		} `json:"response"`
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return
	}
	if data.Response.Status.Code != 0 {
		err = errors.New(data.Response.Status.Message)
		return
	}
	bpm = data.Response.Track.AudioSummary.Tempo
	return
}
