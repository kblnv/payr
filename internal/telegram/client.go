package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Client struct {
	Token string
}

func New(token string) *Client {
	return &Client{Token: token}
}

type sendMessageRequest struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

func (c *Client) SendMessage(chatID int64, text string) error {
	url := fmt.Sprintf(
		"https://api.telegram.org/bot%s/sendMessage",
		c.Token,
	)

	payload := sendMessageRequest{
		ChatID: chatID,
		Text:   text,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	log.Printf("sending message to chat_id=%d", chatID)

	resp, err := http.Post(
		url,
		"application/json",
		bytes.NewBuffer(body),
	)

	if err != nil {
		return fmt.Errorf("http request failed: %w", err)
	}

	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	log.Printf("telegram response status=%d body=%s", resp.StatusCode, string(respBody))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram error: status=%d body=%s", resp.StatusCode, string(respBody))
	}

	log.Println("message successfully sent via telegram")

	return nil
}
