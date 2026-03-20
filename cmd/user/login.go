package user

import (
	"encoding/json"
	"fmt"
	"os"

	"barong-cli/internal/session"
	pkguser "barong-cli/pkg/user"

	"github.com/spf13/cobra"
)

func newLoginCmd(getBaseURL func() string) *cobra.Command {
	var email, password, otpCode string

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Start a new session",
		RunE: func(cmd *cobra.Command, args []string) error {
			client := pkguser.NewClient(getBaseURL())
			u, cookies, err := client.Login(email, password, otpCode)
			if err != nil {
				return err
			}
			if err := session.Save(cookies); err != nil {
				return fmt.Errorf("failed to save session: %w", err)
			}
			fmt.Fprintf(os.Stderr, "Logged in as %s\n", u.Email)
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(u)
		},
	}

	cmd.Flags().StringVar(&email, "email", "", "Account email (required)")
	cmd.Flags().StringVar(&password, "password", "", "Account password (required)")
	cmd.Flags().StringVar(&otpCode, "otp-code", "", "Code from Google Authenticator")
	_ = cmd.MarkFlagRequired("email")
	_ = cmd.MarkFlagRequired("password")

	return cmd
}
