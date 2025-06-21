package main

import (
	"fmt"
	"log"
	"net/http"
	"stream_app/internal/db"
	"stream_app/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error with .env file")
	}

	db.InitPostgres()

	r := chi.NewRouter()

	r.Route("/gameday", func(r chi.Router) {
		r.Post("/new", handlers.CreateGameDay)
		r.Get("/list", handlers.GetAllGameDays)
		r.Delete("/remove", handlers.DeleteDay)
		r.Patch("/save", handlers.SaveDay)
	})

	r.Route("/game", func(r chi.Router) {
		r.Post("/new", handlers.CreateGame)
		r.Delete("/remove", handlers.DeleteGame)
		r.Get("/list", handlers.GetAllGames)
	})

	r.Route("/stats", func(r chi.Router) {
		r.Get("/mood", handlers.GetMoodToday)
		r.Get("/points", handlers.GetPointsToday)
	})

	r.Route("/todaygames", func(r chi.Router) {
		r.Get("/goodminus", handlers.GetGoodMinusGames)
		r.Get("/goodplus", handlers.GetGoodPlusGames)
		r.Get("/badminus", handlers.GetBadMinusGames)
		r.Get("/badplus", handlers.GetBadPlusGames)
	})

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", r)
}
