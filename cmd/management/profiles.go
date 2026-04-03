package management

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
)

func newProfilesCmd(getBaseURL func() string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profiles",
		Short: "Profile management operations",
	}
	cmd.AddCommand(newProfileImportCmd(getBaseURL))
	return cmd
}

func newProfileImportCmd(getBaseURL func() string) *cobra.Command {
	var uid, firstName, lastName, dob, address, postcode, city, country, state, metadata string

	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import a profile for a user",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newManagementClient(getBaseURL())
			if err != nil {
				return err
			}
			u, err := client.ImportProfile(uid, firstName, lastName, dob, address, postcode, city, country, state, metadata)
			if err != nil {
				return err
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(u)
		},
	}
	cmd.Flags().StringVar(&uid, "uid", "", "User UID (required)")
	cmd.Flags().StringVar(&firstName, "first-name", "", "First name")
	cmd.Flags().StringVar(&lastName, "last-name", "", "Last name")
	cmd.Flags().StringVar(&dob, "dob", "", "Date of birth")
	cmd.Flags().StringVar(&address, "address", "", "Address")
	cmd.Flags().StringVar(&postcode, "postcode", "", "Postcode")
	cmd.Flags().StringVar(&city, "city", "", "City")
	cmd.Flags().StringVar(&country, "country", "", "Country")
	cmd.Flags().StringVar(&state, "state", "", "State")
	cmd.Flags().StringVar(&metadata, "metadata", "", "Metadata (JSON string)")
	_ = cmd.MarkFlagRequired("uid")
	return cmd
}
