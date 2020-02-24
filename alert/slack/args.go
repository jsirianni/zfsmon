package slack

import (
    "fmt"


    "github.com/asaskevich/govalidator"
    "github.com/pkg/errors"
)

func (slack Slack) validateArgs(message string) error {
    if slackDebug {
        fmt.Println("validating slack configuration")
        fmt.Println("slack webhook url:", slack.HookURL)
        fmt.Println("slack channel:", slack.Channel)
        fmt.Println("slack message:", message)
    }

	if slack.HookURL == "" {
		return errors.New("slack webhook url is blank")
	}

    if govalidator.IsURL(slack.HookURL) == false {
        return errors.New("slack webhook url is not a valid url")
    }

	if slack.Channel == "" {
		return errors.New("slack channel is blank")
	}

	if message == "" {
		return errors.New("message is blank")
	}

    if slackDebug {
        fmt.Println("slack configuration validation passed")
    }

	return nil
}
