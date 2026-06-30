package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"payr/internal/logger"
	"payr/internal/transports"
)

type Config struct {
	BotToken string `json:"bot_token"`
}

type Telegram struct {
	botToken string
	log      *logger.Logger
}

func New(log *logger.Logger, rawConfig json.RawMessage) (transports.Transport, error) {
	var config Config

	err := json.Unmarshal(rawConfig, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &Telegram{
		botToken: config.BotToken,
		log:      log,
	}, nil
}

type sendMessageRequest struct {
	ChatId int64  `json:"chat_id"`
	Text   string `json:"text"`
}

func (c *Telegram) Send(text string, to string) error {
	chatId, err := strconv.ParseInt(to, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse chat_id: %w", err)
	}

	url := fmt.Sprintf(
		"https://api.telegram.org/bot%v/sendMessage",
		c.botToken,
	)

	payload := sendMessageRequest{
		ChatId: chatId,
		Text:   text,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	c.log.Info("sending message to id=%v", chatId)

	resp, err := http.Post(
		url,
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		c.log.Error("failed to send request: %v", err)
		return err
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.log.Error("failed to read response: %v", err)
		return err
	}

	c.log.Debug("telegram response status=%v body=%v", resp.StatusCode, string(respBody))

	if resp.StatusCode != http.StatusOK {
		c.log.Error("telegram API error: status=%v", resp.StatusCode)
		return fmt.Errorf("telegram API error: status %v", resp.StatusCode)
	}

	return nil
}
