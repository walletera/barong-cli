package user

import "github.com/spf13/cobra"

func NewUserCmd(getBaseURL func() string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Interact with Barong User API",
	}

	cmd.AddCommand(newCreateCmd(getBaseURL))
	cmd.AddCommand(newLoginCmd(getBaseURL))
	cmd.AddCommand(newLogoutCmd(getBaseURL))
	cmd.AddCommand(newMeCmd(getBaseURL))
	cmd.AddCommand(newOTPCmd(getBaseURL))
	cmd.AddCommand(newAPIKeyCmd(getBaseURL))
	cmd.AddCommand(newServiceAccountCmd(getBaseURL))

	return cmd
}
