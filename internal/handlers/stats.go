package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"stream_app/internal/db"
	"time"
)

type GetMoodTodayResponse struct {
	Mood       []int       `json:"mood"`
	TimeStamps []time.Time `json:"timestamps"`
}

func GetMood(w http.ResponseWriter, r *http.Request) {
	dayIdStr := r.URL.Query().Get("dayId")
	dayID, err := strconv.Atoi(dayIdStr)
	if err != nil || dayID <= 0 {
		http.Error(w, "Invalid dayId", http.StatusBadRequest)
		return
	}

	rows, err := db.PG.Query(`
        SELECT mood_diff, game_date
        FROM games
		WHERE day_id = $1
		ORDER BY game_date ASC
    `, dayID)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	row := db.PG.QueryRow(`
        SELECT mood_start
        FROM days
		WHERE id = $1
    `, dayID)

	var moodStart int

	err = row.Scan(&moodStart)
	if err != nil {
		http.Error(w, `{"error": "db scan error"}`, http.StatusInternalServerError)
		return
	}

	mood := []int{moodStart}
	timestamps := []time.Time{time.Time{}}

	for rows.Next() {
		var deltaMood int
		var timestamp time.Time
		if err := rows.Scan(&deltaMood, &timestamp); err != nil {
			http.Error(w, "Scan error", http.StatusInternalServerError)
			return
		}
		mood = append(mood, mood[len(mood)-1]+deltaMood)
		timestamps = append(timestamps, timestamp)
	}

	resp := GetMoodTodayResponse{
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

type GetPointsTodayResponse struct {
	Points int `json:"points"`
}

func GetPoints(w http.ResponseWriter, r *http.Request) {
	dayIdStr := r.URL.Query().Get("dayId")
	dayID, err := strconv.Atoi(dayIdStr)
	if err != nil || dayID <= 0 {
		http.Error(w, "Invalid dayId", http.StatusBadRequest)
		return
	}

	rows, err := db.PG.Query(`
        SELECT points
        FROM games
		WHERE day_id = $1
		ORDER BY game_date ASC
    `, dayID)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	points := 0

	for rows.Next() {
		var p int
		if err := rows.Scan(&p); err != nil {
			http.Error(w, "Scan error", http.StatusInternalServerError)
			return
		}
		points += p
	}

	resp := GetPointsTodayResponse{
		Points: points,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Marshall server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

type GetSomeGamesResponse struct {
	Count int `json:"count"`
}

func GetGoodMinusGames(w http.ResponseWriter, r *http.Request) {
	dayIdStr := r.URL.Query().Get("dayId")
	dayID, err := strconv.Atoi(dayIdStr)
	if err != nil || dayID <= 0 {
		http.Error(w, "Invalid dayId", http.StatusBadRequest)
		return
	}

	row := db.PG.QueryRow(`
        SELECT COUNT(*)
        FROM games
		WHERE day_id = $1 AND status_id = $2
    `, dayID, 1)

	var c int

	err = row.Scan(&c)
	if err != nil {
		http.Error(w, `{"error": "db scan error"}`, http.StatusInternalServerError)
		return
	}

	resp := GetSomeGamesResponse{
		Count: c,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Marshall server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func GetGoodPlusGames(w http.ResponseWriter, r *http.Request) {
	dayIdStr := r.URL.Query().Get("dayId")
	dayID, err := strconv.Atoi(dayIdStr)
	if err != nil || dayID <= 0 {
		http.Error(w, "Invalid dayId", http.StatusBadRequest)
		return
	}

	row := db.PG.QueryRow(`
        SELECT COUNT(*)
        FROM games
		WHERE day_id = $1 AND status_id = $2
    `, dayID, 2)

	var c int

	err = row.Scan(&c)
	if err != nil {
		http.Error(w, `{"error": "db scan error"}`, http.StatusInternalServerError)
		return
	}

	resp := GetSomeGamesResponse{
		Count: c,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Marshall server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func GetBadMinusGames(w http.ResponseWriter, r *http.Request) {
	dayIdStr := r.URL.Query().Get("dayId")
	dayID, err := strconv.Atoi(dayIdStr)
	if err != nil || dayID <= 0 {
		http.Error(w, "Invalid dayId", http.StatusBadRequest)
		return
	}

	row := db.PG.QueryRow(`
        SELECT COUNT(*)
        FROM games
		WHERE day_id = $1 AND status_id = $2
    `, dayID, 3)

	var c int

	err = row.Scan(&c)
	if err != nil {
		http.Error(w, `{"error": "db scan error"}`, http.StatusInternalServerError)
		return
	}

	resp := GetSomeGamesResponse{
		Count: c,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Marshall server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func GetBadPlusGames(w http.ResponseWriter, r *http.Request) {
	dayIdStr := r.URL.Query().Get("dayId")
	dayID, err := strconv.Atoi(dayIdStr)
	if err != nil || dayID <= 0 {
		http.Error(w, "Invalid dayId", http.StatusBadRequest)
		return
	}

	row := db.PG.QueryRow(`
        SELECT COUNT(*)
        FROM games
		WHERE day_id = $1 AND status_id = $2
    `, dayID, 4)

	var c int

	err = row.Scan(&c)
	if err != nil {
		http.Error(w, `{"error": "db scan error"}`, http.StatusInternalServerError)
		return
	}

	resp := GetSomeGamesResponse{
		Count: c,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Marshall server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func GetMoodLast(w http.ResponseWriter, r *http.Request) {
	row := db.PG.QueryRow(`
        SELECT id
        FROM days
		ORDER BY id DESC
		LIMIT 1
    `)

	var dayID int

	err := row.Scan(&dayID)
	if err != nil {
		http.Error(w, `{"error": "db scan error"}`, http.StatusInternalServerError)
		return
	}

	params := r.URL.Query()
	params.Add("dayId", fmt.Sprintf("%d", dayID))

	r.URL.RawQuery = params.Encode()

	GetMood(w, r)
}

func GetPointsLast(w http.ResponseWriter, r *http.Request) {
	row := db.PG.QueryRow(`
        SELECT id
        FROM days
		ORDER BY id DESC
		LIMIT 1
    `)

	var dayID int

	err := row.Scan(&dayID)
	if err != nil {
		http.Error(w, `{"error": "db scan error"}`, http.StatusInternalServerError)
		return
	}

	params := r.URL.Query()
	params.Add("dayId", fmt.Sprintf("%d", dayID))

	r.URL.RawQuery = params.Encode()

	GetPoints(w, r)
}
