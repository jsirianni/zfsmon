package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/jsirianni/zfsmon/zfs"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"
)

var hookURL string
var slackChannel string
var stateFile string
var noAlert bool
var jsonFmt bool

var z zfs.Zfs

var rootCmd = &cobra.Command{
	Use:   "zfsmon",
	Short: "zfs monitoring daemon",
	Run: func(cmd *cobra.Command, args []string) {
		if err := z.ZFSMon(); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&slackChannel, "channel", "", "slack channel")
	rootCmd.PersistentFlags().StringVar(&hookURL, "url", "", "hook url")
	rootCmd.PersistentFlags().StringVar(&stateFile, "state-file", "/tmp/zfsmon", "path for the state file")
	rootCmd.PersistentFlags().BoolVar(&noAlert, "no-alert", false, "do not send alerts")
	rootCmd.PersistentFlags().BoolVar(&jsonFmt, "json", false, "enable json output")
}

func initConfig() {
	if err := checkFlags(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	z = zfs.Zfs{}
	z.HookURL = hookURL
	z.SlackChannel = slackChannel
	z.NoAlert = noAlert
	z.State.File = stateFile
	z.JSONOutput = jsonFmt

}

func checkFlags() error {
	var e error

	if noAlert == true {
		return e
	}

	if len(slackChannel) == 0 {
		e = multierror.Append(e, errors.New("You must pass a channel '--channel <channel_name>' unless '--no-alert' is specified"))
	}

	if len(hookURL) == 0 {
		e = multierror.Append(e, errors.New("You must pass a slack hook url '--url <hook url>' unless --no-alert' is specified"))
	}

	return e
}
