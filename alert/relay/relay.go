package relay

import (
    "net/http"
    "encoding/json"
    "bytes"
    "strconv"
    "io/ioutil"
    "time"

    "github.com/pkg/errors"
)

const DefaultBaseURL = "https://relay.teamitgr.com"
const apiKeyHeader   = "x-relay-api-key"

type Relay struct {
    BaseURL  string
    APIKey   string
}

type payload struct {
    Text string `json:"text"`
}

func (relay Relay) Message(message string) error {
    if relay.BaseURL == "" {
        relay.BaseURL = DefaultBaseURL
    }

    if relay.APIKey == "" {
        return errors.New("relay API Key is not set!")
    }

    b, err := json.Marshal(payload{Text:message})
    if err != nil {
        return err
    }

    req, err := http.NewRequest("POST", relay.BaseURL + "/message", bytes.NewBuffer(b))
    if err != nil {
        return err
    }
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set(apiKeyHeader, relay.APIKey)

    return send(req)
}

func send(req *http.Request) error {
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil
    }
    defer resp.Body.Close()
    if resp.StatusCode != 200 {
        body, _ := ioutil.ReadAll(resp.Body)
        if body == nil {
            body = []byte("")
        }
        return errors.New("Slack returned status: " + strconv.Itoa(resp.StatusCode) + " " + string(body))
    }
    return nil
}
