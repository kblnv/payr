package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"payr/internal/telegram"
)

func main() {
	log.Println("starting payrem...")

	token := os.Getenv("TG_BOT_TOKEN")
	if token == "" {
		log.Fatal("missing TG_BOT_TOKEN")
	}

	chatIdStr := os.Getenv("TG_GROUP_ID")
	if chatIdStr == "" {
		log.Fatal("missing TG_GROUP_ID")
	}

	chatId, err := strconv.ParseInt(chatIdStr, 10, 64)
	if err != nil {
		log.Fatalf("invalid TG_GROUP_ID: %v", err)
	}

	log.Println("config loaded")

	client := telegram.New(token)

	log.Println("starting payr server...")

	mux := http.NewServeMux()

	mux.HandleFunc("/event", func(w http.ResponseWriter, r *http.Request) {
		err = client.SendMessage(chatId, "hello msg")

		if err != nil {
			log.Fatalf("send message failed: %v", err)
		}

		log.Println("message sent successfully")

		w.Write([]byte("ok"))
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("listening on :8080...")

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
