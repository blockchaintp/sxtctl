package commands

import (
	"io"
	"os"

	"github.com/spf13/cobra"
)

var RemoteName string
var Url string
var AccessToken string
var OutputFormat string

func NewCmdRoot(out, errOut io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "sextantctl",
		Short:        "CLI client for sextant.",
		SilenceUsage: true,
		Long: `
sextantctl
==========

Manage your sextant installation.`,
	}

	cmd.AddCommand(NewCmdRemote(out, errOut))
	cmd.AddCommand(NewCmdCluster(out, errOut))
	cmd.AddCommand(NewCmdDeployment(out, errOut))
	cmd.PersistentFlags().StringVarP(&RemoteName, "remote", "r", os.Getenv("SEXTANT_REMOTE"), "The name of the remote sextant api")
	cmd.PersistentFlags().StringVarP(&Url, "url", "u", os.Getenv("SEXTANT_URL"), "The URL of the remote sextant api")
	cmd.PersistentFlags().StringVarP(&AccessToken, "token", "t", os.Getenv("SEXTANT_TOKEN"), "Your testfaster access token")
	cmd.PersistentFlags().StringVarP(&OutputFormat, "output", "o", "text", "The format of the output (text or json)")
	return cmd
}
