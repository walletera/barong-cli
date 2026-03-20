package user

import (
	"fmt"

	"barong-cli/internal/session"
	pkguser "barong-cli/pkg/user"

	"github.com/spf13/cobra"
)

func newLogoutCmd(getBaseURL func() string) *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Destroy current session",
		RunE: func(cmd *cobra.Command, args []string) error {
			cookies, err := session.Load()
			if err != nil {
				return err
			}
			client := pkguser.NewClient(getBaseURL())
			if err := client.Logout(cookies); err != nil {
				return err
			}
			if err := session.Delete(); err != nil {
				return fmt.Errorf("failed to clear local session: %w", err)
			}
			fmt.Println("Logged out")
			return nil
		},
	}
}
