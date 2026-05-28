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
	"payr/internal/transports"
)

type Config struct {
	BotToken string `json:"bot_token"`
	ChatId   string `json:"chat_id"`
}

type Telegram struct {
	botToken string
	chatId   int64
}

func New(rawConfig json.RawMessage) transports.Transport {
	var config Config

	err := json.Unmarshal(rawConfig, &config)
	helpers.Die(err)

	chatId, err := strconv.ParseInt(config.ChatId, 10, 64)
	helpers.Die(err)

	return &Telegram{
		botToken: config.BotToken,
		chatId:   chatId,
	}
}

type sendMessageRequest struct {
	ChatId int64  `json:"chat_id"`
	Text   string `json:"text"`
}

func (c *Telegram) Send(text string) error {
	url := fmt.Sprintf(
		"https://api.telegram.org/bot%v/sendMessage",
		c.botToken,
	)

	payload := sendMessageRequest{
		ChatId: c.chatId,
		Text:   text,
	}

	body, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	log.Printf("sending message to id=%v", c.chatId)

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

func init() {
	transports.RegisterConstructor(
		"telegram",
		New,
	)
}
