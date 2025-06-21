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

type CreateGameRequest struct {
	Place        int `json:"place"`
	Kills        int `json:"kills"`
	Liquidations int `json:"liquidations"`
	Damage       int `json:"damage"`
	Mood_diff    int `json:"mood_diff"`
	Points       int `json:"points"`
	StatusCode   int `json:"status_code"`
}

type CreateGameResponse struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}

func CreateGame(w http.ResponseWriter, r *http.Request) {
	var dayID int
	dayIdStr := r.URL.Query().Get("dayId")
	if dayIdStr != "" {
		var err error
		dayID, err = strconv.Atoi(dayIdStr)
		if err != nil || dayID <= 0 {
			http.Error(w, "Invalid dayId", http.StatusBadRequest)
			return
		}
	}

	tx, err := db.PG.Begin()
	if err != nil {
		http.Error(w, `{"error": "db error"}`, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	if dayID == 0 {
		var day models.GameDay
		row := tx.QueryRow(`
			SELECT id, saved
			FROM days 
			ORDER BY playday DESC
			LIMIT 1
		`)

		err = row.Scan(&day.ID, &day.Saved)
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
		dayID = day.ID
	}

	var saved bool
	row := tx.QueryRow(`
		SELECT saved
		FROM days 
		WHERE id = $1
	`, dayID)

	err = row.Scan(&saved)
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
	if saved {
		http.Error(w, `{"error": "this day already saved"}`, http.StatusBadRequest)
		return
	}

	var req CreateGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var id int
	var createdAt time.Time
	err = tx.QueryRow(`
		INSERT INTO games (place, kills, liquidations, damage, mood_diff, points, status_id, day_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, game_date
	`, req.Place, req.Kills, req.Liquidations, req.Damage, req.Mood_diff, req.Points, req.StatusCode, dayID).Scan(&id, &createdAt)

	if err != nil {
		http.Error(w, "Insert failed", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Transaction failed", http.StatusInternalServerError)
		return
	}

	resp := CreateGameResponse{
		ID:        id,
		CreatedAt: createdAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type GetAllGamesResponse struct {
	Games []models.Game `json:"games"`
}

func GetAllGames(w http.ResponseWriter, r *http.Request) {
	dayIdStr := r.URL.Query().Get("dayId")
	dayID, err := strconv.Atoi(dayIdStr)
	if err != nil || dayID <= 0 {
		http.Error(w, "Invalid dayId", http.StatusBadRequest)
		return
	}

	rows, err := db.PG.Query(`
        SELECT *
        FROM games
		WHERE day_id = $1
		ORDER BY game_date ASC
    `, dayID)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var games []models.Game
	for rows.Next() {
		var game models.Game
		if err := rows.Scan(&game.ID, &game.Place, &game.Kills, &game.Liquidations, &game.Damage, &game.Mood_diff, &game.Points, &game.StatusCode, &game.GameDate, &game.DayID); err != nil {
			http.Error(w, "Scan error", http.StatusInternalServerError)
			return
		}
		games = append(games, game)
	}

	resp := GetAllGamesResponse{
		Games: games,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Marshall server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func DeleteGame(w http.ResponseWriter, r *http.Request) {
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

	var game models.Game
	row := tx.QueryRow(`
		SELECT *
		FROM games WHERE id = $1 FOR UPDATE
	`, id)

	err = row.Scan(&game.ID, &game.Place, &game.Kills, &game.Liquidations, &game.Damage, &game.Mood_diff, &game.Points, &game.StatusCode, &game.GameDate, &game.DayID)
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
		DELETE FROM games WHERE id = $1
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
		"id":      game.ID,
		"gameday": game.GameDate,
	})
}
