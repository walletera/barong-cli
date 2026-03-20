package user

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
)

func newMeCmd(getBaseURL func() string) *cobra.Command {
	return &cobra.Command{
		Use:   "me",
		Short: "Return current user information",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newAuthenticatedClient(getBaseURL())
			if err != nil {
				return err
			}
			u, err := client.GetMe()
			if err != nil {
				return err
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(u)
		},
	}
}
