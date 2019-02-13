package cmd
import (
    "os"
    "fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "zfsmon",
	Short: "zfs monitoring daemon",
    Run: func(cmd *cobra.Command, args []string) {
        err := run()
        if err != nil {
            fmt.Println(err.Error())
            os.Exit(1)
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
	rootCmd.PersistentFlags().StringVar(&channel, "channel", "", "slack channel")
	rootCmd.PersistentFlags().StringVar(&hook_url, "url", "/opt/zfsmon/alerts.dat", "hook url" )
    rootCmd.PersistentFlags().StringVar(&alertFile, "alert-file", "/usr/", "hook url" )
    rootCmd.PersistentFlags().BoolVar(&daemon, "daemon", false, "daemonize")
    rootCmd.PersistentFlags().IntVar(&checkInt, "interval", 5, "how often to check the zpool(s) in  minutes")
    rootCmd.PersistentFlags().BoolVar(&printReport, "print", false, "print the health report")
    rootCmd.PersistentFlags().BoolVar(&noAlert, "no-alert", false, "do not send alerts")
}
