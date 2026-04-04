package user

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
)

func newServiceAccountCmd(getBaseURL func() string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service-account",
		Short: "Manage service accounts for the current user",
	}

	cmd.AddCommand(newServiceAccountListCmd(getBaseURL))

	return cmd
}

func newServiceAccountListCmd(getBaseURL func() string) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all service accounts for the current user",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newAuthenticatedClient(getBaseURL())
			if err != nil {
				return err
			}
			accounts, err := client.ListServiceAccounts()
			if err != nil {
				return err
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(accounts)
		},
	}
}
