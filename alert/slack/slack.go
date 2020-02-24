package slack

import (
	"os"
	"fmt"
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
)

const envSlackDebug = "SLACK_DEBUG"
var slackDebug = false

type Slack struct {
	HookURL string
	Channel string
}

type Payload struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

func (slack Slack) Print() {
	fmt.Println("slack hook url:", slack.HookURL)
	fmt.Println("slack channel:", slack.Channel)
}

func (slack Slack) Message(message string) error {
	// set debug, ignore parse errors
	x := os.Getenv(envSlackDebug)
	slackDebug, _ = strconv.ParseBool(x)

	if err := slack.validateArgs(message); err != nil {
		return errors.Wrap(err, "slack configuration failed validation")
	}

	return slack.sendPayload(message)
}

func (slack Slack) sendPayload(m string) error {
	payload := Payload{
		Channel: slack.Channel,
		Text:    m,
	}

	p, err := json.Marshal(payload)
	if err != nil {
		return nil
	}

	if slackDebug {
		fmt.Println("slack payload: " + string(p))
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
