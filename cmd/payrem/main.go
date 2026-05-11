package main

import (
	"log"
	"os"
	"strconv"

	"payrem/internal/telegram"
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

	err = client.SendMessage(chatId, "hello msg")
	if err != nil {
		log.Fatalf("send message failed: %v", err)
	}

	log.Println("message sent successfully")
}