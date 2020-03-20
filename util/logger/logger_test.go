package logger

import (
    "testing"
)

func TestConfigure(t *testing.T) {
    levels := []string{
        "trace",
        "info",
        "warning",
        "error",
    }

    l := Logger{}

    for _, level := range levels {
        if err := l.Configure(level); err != nil {
            t.Error(err)
            continue
        }

        if l.logLevel != level {
            t.Errorf("expected l.logLevl to be set to " + level + " got " + l.logLevel)
        }

        if l.LogLevel() != level {
            t.Errorf("expected l.LogLevel() to return " + level + " got " + l.LogLevel())
        }

        if l.configured != true {
            t.Errorf("expected l.configured to be true, got false")
        }

        if l.Configured() != true {
            t.Errorf("expected l.Configured() to return true, got false")
        }
    }
}
