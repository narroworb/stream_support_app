package handlers

import (
	"encoding/json"
	"net/http"
	"stream_app/internal/db"
	"time"
)

type GetMoodHistoryResponse struct {
	Mood       []int       `json:"mood"`
	TimeStamps []time.Time `json:"timestamps"`
}

func GetMoodHistory(w http.ResponseWriter, r *http.Request) {
	rows, err := db.PG.Query(`
        SELECT mood_end - mood_start, playday
        FROM days
		ORDER BY playday ASC
    `)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	mood := []int{}
	timestamps := []time.Time{}

	for rows.Next() {
		var deltaMood int
		var timestamp time.Time
		if err := rows.Scan(&deltaMood, &timestamp); err != nil {
			http.Error(w, "Scan error", http.StatusInternalServerError)
			return
		}
		mood = append(mood, deltaMood)
		timestamps = append(timestamps, timestamp)
	}

	resp := GetMoodHistoryResponse{
		Mood:       mood,
		TimeStamps: timestamps,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Marshall server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

type GetPointsHistoryResponse struct {
	Points     []int         `json:"points"`
	TimeStamps []time.Time `json:"timestamps"`
}

func GetPointsHistory(w http.ResponseWriter, r *http.Request) {
	rows, err := db.PG.Query(`
        SELECT total_points, playday
        FROM days
		ORDER BY playday ASC
    `)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	points := []int{}
	timestamps := []time.Time{}

	for rows.Next() {
		var p int
		var timestamp time.Time
		if err := rows.Scan(&p, &timestamp); err != nil {
			http.Error(w, "Scan error", http.StatusInternalServerError)
			return
		}
		points = append(points, p)
		timestamps = append(timestamps, timestamp)
	}

	resp := GetPointsHistoryResponse{
		Points:       points,
		TimeStamps: timestamps,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Marshall server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
