package relay

import (
    "testing"
)

func TestMessagE(t *testing.T) {
    r := Relay{}
    if r.Message("") == nil {
        t.Errorf("relay.Message() should return an error when given an empty message")
    }
}

func TestInit(t *testing.T) {
    // uses a random uuid
    r := Relay{
        BaseURL: "",
        APIKey: "93733f41-2952-4ecc-a36c-c32782ca5ce5",
    }

    if err := r.init(); err != nil {
        t.Errorf(err.Error())
        return
    }

    if r.BaseURL != DefaultBaseURL {
        t.Errorf("expected base url to be set by init when empty")
    }
}
