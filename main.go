package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "modernc.org/sqlite" // Pure Go driver
)

type Review struct {
	ID           string `json:"review_id"`
	LocationID   string `json:"location_id"`
	Author       string `json:"author"`
	ProfilePhoto string `json:"profile_photo"`
	Rating       string `json:"rating"`
	Text         string `json:"text"`
	MediaLinks   string `json:"media_links"`
}

var db *sql.DB

func initDB() {
	var err error
	// Use "sqlite" instead of "sqlite3" for the Pure Go driver
	db, err = sql.Open("sqlite", "./reviews.db")
	if err != nil {
		log.Fatal(err)
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS reviews (
		"id" TEXT PRIMARY KEY,
		"location_id" TEXT,
		"reviewer_name" TEXT,
		"profile_photo" TEXT,
		"star_rating" TEXT,
		"comment" TEXT,
		"media_urls" TEXT,
		"created_at" DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal("Table creation failed:", err)
	}
	fmt.Println("SQLite Database (Pure Go) initialized successfully.")
}

func reviewHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		return
	}

	if r.Method == "POST" {
		var reviews []Review
		if err := json.NewDecoder(r.Body).Decode(&reviews); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		for _, rev := range reviews {
			insertSQL := `INSERT INTO reviews (id, location_id, reviewer_name, profile_photo, star_rating, comment, media_urls) 
						  VALUES (?, ?, ?, ?, ?, ?, ?)
						  ON CONFLICT(id) DO UPDATE SET
						  star_rating=excluded.star_rating,
						  comment=excluded.comment;`
			
			_, err := db.Exec(insertSQL, rev.ID, rev.LocationID, rev.Author, rev.ProfilePhoto, rev.Rating, rev.Text, rev.MediaLinks)
			if err != nil {
				log.Printf("Could not save review %s: %v", rev.ID, err)
			}
		}
		fmt.Printf("Received and saved %d reviews.\n", len(reviews))
		w.WriteHeader(http.StatusCreated)
	}
}

func main() {
	initDB()
	defer db.Close()
	http.HandleFunc("/api/reviews", reviewHandler)
	fmt.Println("Server starting at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}