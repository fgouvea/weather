package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

type AppConfig struct {
	Port string
}

func readConfigFromEnv() AppConfig {
	return AppConfig{
		Port: fmt.Sprintf(":%s", readFromEnv("PORT", "8083")),
	}
}

type Notification struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

func main() {
	config := readConfigFromEnv()

	r := chi.NewRouter()

	r.Route("/notification/send", func(r chi.Router) {
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			var body Notification
			err := json.NewDecoder(r.Body).Decode(&body)

			if err != nil {
				fmt.Println("ERROR READING NOTIFICATION")
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			fmt.Println("NOTIFICATION RECEIVED")
			fmt.Printf("ID: %s\n", body.ID)
			fmt.Println(body.Content)

			w.WriteHeader(http.StatusAccepted)
		})
	})

	fmt.Printf("Listenning on %s...\n", config.Port)
	http.ListenAndServe(config.Port, r)
}

func readFromEnv(env, def string) string {
	if value := os.Getenv(env); value != "" {
		return value
	}

	return def
}
