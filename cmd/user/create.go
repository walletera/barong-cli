package user

import (
	"encoding/json"
	"os"

	pkguser "barong-cli/pkg/user"

	"github.com/spf13/cobra"
)

func newCreateCmd(getBaseURL func() string) *cobra.Command {
	var email, password, username, refid string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new user",
		RunE: func(cmd *cobra.Command, args []string) error {
			client := pkguser.NewClient(getBaseURL())
			u, err := client.CreateUser(email, password, username, refid)
			if err != nil {
				return err
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(u)
		},
	}

	cmd.Flags().StringVar(&email, "email", "", "User email (required)")
	cmd.Flags().StringVar(&password, "password", "", "User password (required)")
	cmd.Flags().StringVar(&username, "username", "", "Username")
	cmd.Flags().StringVar(&refid, "refid", "", "Referral UID")
	_ = cmd.MarkFlagRequired("email")
	_ = cmd.MarkFlagRequired("password")

	return cmd
}
