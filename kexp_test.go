package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func expectedNil(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("expected nil value but got:%v", err)
	}
}

func expectedDeepEqual(t *testing.T, v, k interface{}) {
	if !reflect.DeepEqual(v, k) {
		t.Fatalf("expected these to be equal:\nACTUAL:\n%v\nEXPECTED:\n%v", v, k)
	}
}

func TestNowPlaying(t *testing.T) {
	for _, tt := range []struct {
		kexp     Kexp
		response string
	}{
		{
			response: "...air break...",
			kexp: Kexp{Plays: []Play{
				{
					Track: Information{Name: ""},
				},
			}},
		},
		{
			response: "Wandering Star by Portishead from Dummy :: released in 1994",
			kexp: Kexp{
				Plays: []Play{
					{
						Track:   Information{Name: "Wandering Star"},
						Artist:  Information{Name: "Portishead"},
						Release: Information{Name: "Dummy"},
						ReleaseEvent: ReleaseEvent{
							Date: Date{Year: 1994},
						},
					},
				},
			},
		},
	} {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := json.NewEncoder(w).Encode(tt.kexp)
			expectedNil(t, err)
			w.Header().Set("Content-Type", "application/json")
		}))
		defer srv.Close()

		np, err := NowPlaying(srv.URL)
		expectedNil(t, err)
		expectedDeepEqual(t, np[0], tt.response)
	}
}
