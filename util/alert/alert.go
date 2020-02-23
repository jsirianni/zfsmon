package alert

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type Slack struct {
	HookURL string
	Post    SlackPost
}

type SlackPost struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

func (slack *Slack) BasicMessage() error {
	return slack.sendPayload()
}

func (slack *Slack) sendPayload() error {
	payload, err := json.Marshal(slack.Post)
	if err != nil {
		return nil
	}

	req, err := http.NewRequest("POST", slack.HookURL, bytes.NewBuffer(payload))
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
