package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

var (
	kexpURL = "http://cache.kexp.org/cache/plays"
)

// Date describes the date where the album was released.
type Date struct {
	Year  int `json:"Year"`
	Month int `json:"Month"`
	Day   int `json:"Day"`
}

// ReleaseEvent describes the date and the country where the album was released.
type ReleaseEvent struct {
	Type    string `json:"__type"`
	URI     string `json:"Uri"`
	Country string `json:"Country"`
	Date    Date   `json:"Date"`
}

// Information describes common informations of Artist, Release and Track playing on KEXP.
type Information struct {
	Type string `json:"__type"`
	URI  string `json:"Uri"`
	ID   string `json:"MusicBrainzId"`
	Name string `json:"Name"`
}

// Play describes the current playlist on KEXP.
type Play struct {
	URI          string       `json:"Uri"`
	Type         int          `json:"Type"`
	RotationType int          `json:"RotationType"`
	IsRequest    bool         `json:"IsRequest"`
	AirDate      string       `json:"AirDate"`
	IsDeleted    bool         `json:"IsDeleted"`
	CreatedDate  string       `json:"CreatedDate"`
	UpdatedDate  string       `json:"UpdatedDate"`
	Artist       Information  `json:"Artist"`
	Release      Information  `json:"Release"`
	Track        Information  `json:"Track"`
	ReleaseEvent ReleaseEvent `json:"ReleaseEvent"`
}

// String returns a formated string with play information.
func (p *Play) String() string {
	if p.Track.Name == "" && p.Artist.Name == "" && p.Release.Name == "" {
		return "...air break..."
	}
	if p.Track.Name == "" {
		return fmt.Sprintf("%s from %s :: released in %d",
			p.Artist.Name, p.Release.Name, p.ReleaseEvent.Date.Year)
	}
	if p.ReleaseEvent.Date.Year == 0 {
		return fmt.Sprintf("%s by %s from %s",
			p.Track.Name, p.Artist.Name, p.Release.Name)
	}
	return fmt.Sprintf("%s by %s from %s :: released in %d",
		p.Track.Name, p.Artist.Name, p.Release.Name,
		p.ReleaseEvent.Date.Year)
}

// Kexp describes what's now playing on KEXP.
type Kexp struct {
	Count        int    `json:"Count"`
	IsStale      bool   `json:"IsStale"`
	PlayDate     string `json:"PlayDate"`
	ResponseDate string `json:"ResponseDate"`
	Plays        []Play `json:"Plays"`
}

// NowPlaying returns list of songs that are playing on KEXP.
func NowPlaying(host string) ([]string, error) {
	u, err := url.Parse(host)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	conn := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := conn.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var kexp Kexp
	if err := json.Unmarshal(body, &kexp); err != nil {
		return nil, err
	}

	var np []string
	for _, play := range kexp.Plays {
		np = append(np, play.String())
	}
	return np, nil
}
