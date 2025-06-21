package models

import "time"

type Game struct {
	ID           int       `json:"id"`
	Place        int       `json:"place"`
	Kills        int       `json:"kills"`
	Liquidations int       `json:"liquidations"`
	Damage       int       `json:"damage"`
	Mood_diff    int       `json:"mood_diff"`
	Points       int       `json:"points"`
	StatusCode   int       `json:"status_code"`
	DayID        int       `json:"day_id"`
	GameDate     time.Time `json:"game_date"`
}
