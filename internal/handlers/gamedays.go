package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"stream_app/internal/db"
	"stream_app/internal/models"
	"time"
)

type CreateGameDayResponse struct {
	ID        int       `json:"id"`
	MoodStart int       `json:"mood_start"`
	CreatedAt time.Time `json:"createdAt"`
}

func CreateGameDay(w http.ResponseWriter, r *http.Request) {
	var day models.GameDay
	if err := json.NewDecoder(r.Body).Decode(&day); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if day.MoodStart == 0 {
		http.Error(w, "MoodStart is required", http.StatusBadRequest)
		return
	}

	tx, err := db.PG.Begin()
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	var id int
	var createdAt time.Time
	err = tx.QueryRow(`
		INSERT INTO days (mood_start)
		VALUES ($1)
		RETURNING id, playday
	`, day.MoodStart).Scan(&id, &createdAt)

	if err != nil {
		http.Error(w, "Insert failed", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Transaction failed", http.StatusInternalServerError)
		return
	}

	resp := CreateGameDayResponse{
		ID:        id,
		MoodStart: day.MoodStart,
		CreatedAt: createdAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type GetAllGameDaysResponse struct {
	Days []models.GameDay `json:"gamedays"`
}

func GetAllGameDays(w http.ResponseWriter, r *http.Request) {
	rows, err := db.PG.Query(`
        SELECT *
        FROM days
		ORDER BY playday DESC
    `)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var days []models.GameDay
	for rows.Next() {
		var d models.GameDay
		if err := rows.Scan(&d.ID, &d.MoodStart, &d.MoodEnd, &d.TotalPoints, &d.PlayDay, &d.Saved); err != nil {
			http.Error(w, "Scan error", http.StatusInternalServerError)
			return
		}
		days = append(days, d)
	}

	resp := GetAllGameDaysResponse{
		Days: days,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Marshall server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func DeleteDay(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.Error(w, `{"error": "invalid or missing id"}`, http.StatusBadRequest)
		return
	}

	tx, err := db.PG.Begin()
	if err != nil {
		http.Error(w, `{"error": "db error"}`, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	var day models.GameDay
	row := tx.QueryRow(`
		SELECT *
		FROM days WHERE id = $1 FOR UPDATE
	`, id)

	err = row.Scan(&day.ID, &day.MoodStart, &day.MoodEnd, &day.TotalPoints, &day.PlayDay, &day.Saved)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    3,
			"message": "errors.common.notFound",
			"details": struct{}{},
		})
		return
	} else if err != nil {
		http.Error(w, `{"error": "db scan error"}`, http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(`
		DELETE FROM days WHERE id = $1
	`, id)
	if err != nil {
		http.Error(w, `{"error": "delete failed"}`, http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, `{"error": "commit failed"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      day.ID,
		"gameday": day.PlayDay,
	})
}

func SaveDay(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.Error(w, `{"error": "invalid or missing id"}`, http.StatusBadRequest)
		return
	}

	tx, err := db.PG.Begin()
	if err != nil {
		http.Error(w, `{"error": "db error"}`, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	var day models.GameDay
	row := tx.QueryRow(`
		SELECT *
		FROM days WHERE id = $1 FOR UPDATE
	`, id)

	err = row.Scan(&day.ID, &day.MoodStart, &day.MoodEnd, &day.TotalPoints, &day.PlayDay, &day.Saved)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    3,
			"message": "errors.common.notFound",
			"details": struct{}{},
		})
		return
	} else if err != nil {
		http.Error(w, `{"error": "db scan error"}`, http.StatusInternalServerError)
		return
	}
	if day.Saved {
		http.Error(w, `{"error": "already saved"}`, http.StatusBadRequest)
		return
	}

	rows, err := db.PG.Query(`
        SELECT mood_diff, points
        FROM games
		WHERE day_id = $1;
    `, day.ID)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var pointsCounter int
	var moodCounter int
	for rows.Next() {
		var m, p int
		if err := rows.Scan(&m, &p); err != nil {
			fmt.Println(err)
			http.Error(w, "Scan error", http.StatusInternalServerError)
			return
		}
		pointsCounter += p
		moodCounter += m
	}

	_, err = tx.Exec(`
		UPDATE days SET saved = $1, mood_end = $2, total_points = $3 WHERE id = $4
	`, true, day.MoodStart+moodCounter, pointsCounter, id)
	if err != nil {
		http.Error(w, `{"error": "update failed"}`, http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, `{"error": "commit failed"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":           day.ID,
		"mood_start":   day.MoodStart,
		"mood_end":     day.MoodStart+moodCounter,
		"total_points": pointsCounter,
		"playday":      day.PlayDay,
	})
}
