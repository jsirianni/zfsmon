package slack

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type Slack struct {
	HookURL string
	Channel string
}

type payload struct {
	channel string `json:"channel"`
	text    string `json:"text"`
}

func (slack Slack) Message(message string) error {
	return slack.sendPayload(message)
}

func (slack Slack) sendPayload(m string) error {
	payload := payload{
		channel: slack.Channel,
		text:    m,
	}

	p, err := json.Marshal(payload)
	if err != nil {
		return nil
	}

	req, err := http.NewRequest("POST", slack.HookURL, bytes.NewBuffer(p))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("Slack returned status: " + strconv.Itoa(resp.StatusCode))
	}
	return nil
}
