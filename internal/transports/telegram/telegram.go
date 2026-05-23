package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"payr/internal/helpers"
)

type Telegram struct {
	sender    string
	channelId int64
}

func New(sender string, channelId string) *Telegram {
	parsedChannelId, err := strconv.ParseInt(channelId, 10, 64)
	helpers.Die(err)

	return &Telegram{sender: sender, channelId: parsedChannelId}
}

type sendMessageRequest struct {
	ChatId int64  `json:"chat_id"`
	Text   string `json:"text"`
}

func (c *Telegram) Send(text string) error {
	url := fmt.Sprintf(
		"https://api.telegram.org/bot%v/sendMessage",
		c.sender,
	)

	payload := sendMessageRequest{
		ChatId: c.channelId,
		Text:   text,
	}

	body, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	log.Printf("sending message to id=%v", c.channelId)

	resp, err := http.Post(
		url,
		"application/json",
		bytes.NewBuffer(body),
	)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	log.Printf("telegram response status=%v body=%v", resp.StatusCode, string(respBody))

	return nil
}

func (c *Telegram) Name() string {
	return "telegram"
}
