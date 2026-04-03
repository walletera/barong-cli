package management

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newTimestampCmd(getBaseURL func() string) *cobra.Command {
	return &cobra.Command{
		Use:   "timestamp",
		Short: "Return server time in seconds since Unix epoch",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newManagementClient(getBaseURL())
			if err != nil {
				return err
			}
			ts, err := client.GetTimestamp()
			if err != nil {
				return err
			}
			fmt.Fprintf(os.Stdout, "%d\n", ts)
			return nil
		},
	}
}
