package models

import "time"

type GameDay struct {
	ID          int       `json:"id"`
	MoodStart   int       `json:"mood_start"`
	MoodEnd     int       `json:"mood_end"`
	TotalPoints int       `json:"total_points"`
	PlayDay     time.Time `json:"playday"`
	Saved       bool      `json:"saved"`
}
