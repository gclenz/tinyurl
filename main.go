package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type UrlData struct {
	ID        string    `json:"id"`
	Url       string    `json:"url"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	db := GetDatabaseConnection()
	repo := NewUrlRepository(db)
	controller := NewController(repo)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /urls", controller.CreateUrl)
	mux.HandleFunc("GET /urls/{id}", controller.GetUrl)
	mux.HandleFunc("GET /healthz", controller.Healthz)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), mux)
	if err != nil {
		log.Fatalln(err)
	}
}
