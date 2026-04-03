package management

import (
	"github.com/spf13/cobra"
)

var (
	keyID          string
	privateKeyFile string
)

func NewManagementCmd(getBaseURL func() string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "management",
		Short: "Interact with Barong Management API",
	}

	cmd.PersistentFlags().StringVar(&keyID, "key-id", "", "Management API key ID (overrides BARONG_MANAGEMENT_KEY_ID)")
	cmd.PersistentFlags().StringVar(&privateKeyFile, "private-key-file", "", "Path to RSA private key PEM file (overrides BARONG_MANAGEMENT_PRIVATE_KEY_FILE)")

	cmd.AddCommand(newUsersCmd(getBaseURL))
	cmd.AddCommand(newLabelsCmd(getBaseURL))
	cmd.AddCommand(newProfilesCmd(getBaseURL))
	cmd.AddCommand(newPhonesCmd(getBaseURL))
	cmd.AddCommand(newDocumentsCmd(getBaseURL))
	cmd.AddCommand(newServiceAccountsCmd(getBaseURL))
	cmd.AddCommand(newOTPCmd(getBaseURL))
	cmd.AddCommand(newTimestampCmd(getBaseURL))

	return cmd
}
