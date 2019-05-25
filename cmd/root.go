package cmd
import (
    "os"
    "fmt"
    "errors"

    "zfsmon/zfs"

	"github.com/spf13/cobra"
    multierror "github.com/hashicorp/go-multierror"

)

var hookURL      string
var slackChannel string
var alertFile    string
var noAlert      bool

var z zfs.Zfs

var rootCmd = &cobra.Command{
	Use:   "zfsmon",
	Short: "zfs monitoring daemon",
    Run: func(cmd *cobra.Command, args []string) {
        if err :=  z.ZFSMon(); err != nil {
            fmt.Fprintln(os.Stderr, err.Error())
            os.Exit(1)
        } else {
            os.Exit(0)
        }
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
	rootCmd.PersistentFlags().StringVar(&hookURL, "url", "", "hook url" )
    rootCmd.PersistentFlags().StringVar(&alertFile, "alert-file", "/tmp/zfsmon", "hook url" )
    rootCmd.PersistentFlags().BoolVar(&noAlert, "no-alert", false, "do not send alerts")
}

func initConfig() {
    if err := checkFlags(); err != nil {
       fmt.Fprintln(os.Stderr, err.Error())
       os.Exit(1)
   }

    z = zfs.Zfs{
        HookURL: hookURL,
        SlackChannel: slackChannel,
        NoAlert: noAlert,
        AlertFile: alertFile,
    }
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
