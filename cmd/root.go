package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/jsirianni/zfsmon/alert/slack"
	"github.com/jsirianni/zfsmon/alert/relay"
	"github.com/jsirianni/zfsmon/alert/terminal"
	"github.com/jsirianni/zfsmon/zfs"

	"github.com/spf13/cobra"
)

const defaultLogLevl = "error"

var hookURL string
var slackChannel string
var stateFile string
var alertType string
var noAlert bool
var daemon bool
var logLevel string

var (
	relayHost string
	relayAPIKey string
)

var z zfs.Zfs

var rootCmd = &cobra.Command{
	Use:   "zfsmon",
	Short: "zfs monitoring util",
	Run: func(cmd *cobra.Command, args []string) {
		if err := z.ZFSMon(); err != nil {
			z.Log.Error(err)
			os.Exit(1)
		}
		os.Exit(0)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&stateFile, "state-file", "", "path for the state file")
	rootCmd.PersistentFlags().BoolVar(&daemon, "daemon", false, "enable daemon mode")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", defaultLogLevl, "logging level [error, warning, info, trace]")

	// alert flags
	rootCmd.PersistentFlags().BoolVar(&noAlert, "no-alert", false, "do not send alerts")
	rootCmd.PersistentFlags().StringVar(&alertType, "alert-type", "", "alert system to use")

	// slack alert type
	rootCmd.PersistentFlags().StringVar(&slackChannel, "slack-channel", "", "slack channel")
	rootCmd.PersistentFlags().StringVar(&hookURL, "slack-url", "", "hook url")

	// relay alert type
	rootCmd.PersistentFlags().StringVar(&relayHost, "relay-host", "", "relay host URL")
	rootCmd.PersistentFlags().StringVar(&relayAPIKey, "relay-api-key", "", "relay api key")
}

func initConfig() {
	if err := initFlags(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	if err := z.Init(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func initFlags() error {
	if err := z.Log.Configure(logLevel); err != nil {
		return err
	}

	z.DaemonMode = daemon

	if err := initSate(); err != nil {
		return err
	}
	if err := initAlert(); err != nil {
		return err
	}
	if err := initHostname(); err != nil {
		return err
	}
	return nil
}

func initSate() error {
	if stateFile == "" {
		return errors.New("state file must be configured")
	}
	z.State.File = stateFile
	return nil
}

func initHostname() error {
	var err error
	z.Hostname, err = os.Hostname()
	if err != nil {
		z.Hostname = "could_not_detect_hostname_sorry"
		fmt.Println("could not detect hostname")
	}
	return nil
}

func initAlert() error {
	if noAlert {
		z.AlertConfig.NoAlert = noAlert
		return nil
	}

	if alertType == "" {
		return errors.New("--alert-type is not set")
	}
	if alertType == "slack" {
		return initSlack()
	}
	if alertType == "relay" {
		return initRelay()
	}
	if alertType == "terminal" {
		return initTerminal()
	}
	return errors.New("alert type not valid: " + alertType)
}

func initSlack() error {
	if hookURL == "" {
		return errors.New("hook url must be set when alert type is 'slack'")
	}
	if slackChannel == "" {
		return errors.New("slack channel must be set when alert type is 'slack'")
	}
	z.Alert = slack.Slack{hookURL, slackChannel}
	return nil
}

func initRelay() error {
	if relayAPIKey == "" {
		return errors.New("relay api key must be set when alert type is 'relay'")
	}
	r := relay.Relay{
		BaseURL: relayHost,
		APIKey: relayAPIKey,
	}
	z.Alert = r
	z.Log.Trace(r.APIKey)
	return nil
}

func initTerminal() error {
	z.Alert = terminal.Terminal{}
	return nil
}
