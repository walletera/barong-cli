package management

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newUsersCmd(getBaseURL func() string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "users",
		Short: "User management operations",
	}
	cmd.AddCommand(newUserCreateCmd(getBaseURL))
	cmd.AddCommand(newUserGetCmd(getBaseURL))
	cmd.AddCommand(newUserListCmd(getBaseURL))
	cmd.AddCommand(newUserUpdateCmd(getBaseURL))
	cmd.AddCommand(newUserImportCmd(getBaseURL))
	return cmd
}

func newUserCreateCmd(getBaseURL func() string) *cobra.Command {
	var email, password, referralUID string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new user",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newManagementClient(getBaseURL())
			if err != nil {
				return err
			}
			u, err := client.CreateUser(email, password, referralUID)
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
	cmd.Flags().StringVar(&referralUID, "referral-uid", "", "Referral UID")
	_ = cmd.MarkFlagRequired("email")
	_ = cmd.MarkFlagRequired("password")
	return cmd
}

func newUserGetCmd(getBaseURL func() string) *cobra.Command {
	var uid, email, phoneNum string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get user and profile information",
		RunE: func(cmd *cobra.Command, args []string) error {
			if uid == "" && email == "" && phoneNum == "" {
				return fmt.Errorf("at least one of --uid, --email, or --phone must be provided")
			}
			client, err := newManagementClient(getBaseURL())
			if err != nil {
				return err
			}
			u, err := client.GetUser(uid, email, phoneNum)
			if err != nil {
				return err
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(u)
		},
	}
	cmd.Flags().StringVar(&uid, "uid", "", "User UID")
	cmd.Flags().StringVar(&email, "email", "", "User email")
	cmd.Flags().StringVar(&phoneNum, "phone", "", "User phone number")
	return cmd
}

func newUserListCmd(getBaseURL func() string) *cobra.Command {
	var extended bool
	var from, to, page, limit int64

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List users",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newManagementClient(getBaseURL())
			if err != nil {
				return err
			}
			users, err := client.ListUsers(extended, from, to, page, limit)
			if err != nil {
				return err
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(users)
		},
	}
	cmd.Flags().BoolVar(&extended, "extended", false, "Return full user information")
	cmd.Flags().Int64Var(&from, "from", 0, "Unix timestamp: return records from this time")
	cmd.Flags().Int64Var(&to, "to", 0, "Unix timestamp: return records before this time")
	cmd.Flags().Int64Var(&page, "page", 0, "Page number (default 1)")
	cmd.Flags().Int64Var(&limit, "limit", 0, "Users per page (max 100)")
	return cmd
}

func newUserUpdateCmd(getBaseURL func() string) *cobra.Command {
	var uid, role, data string

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update role and data fields of an existing user",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newManagementClient(getBaseURL())
			if err != nil {
				return err
			}
			u, err := client.UpdateUser(uid, role, data)
			if err != nil {
				return err
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(u)
		},
	}
	cmd.Flags().StringVar(&uid, "uid", "", "User UID (required)")
	cmd.Flags().StringVar(&role, "role", "", "New user role")
	cmd.Flags().StringVar(&data, "data", "", "Additional key:value pairs in JSON format")
	_ = cmd.MarkFlagRequired("uid")
	return cmd
}

func newUserImportCmd(getBaseURL func() string) *cobra.Command {
	var email, passwordDigest, referralUID, phone string
	var firstName, lastName, dob, address, postcode, city, country, state string

	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import an existing user",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newManagementClient(getBaseURL())
			if err != nil {
				return err
			}
			u, err := client.ImportUser(email, passwordDigest, referralUID, phone, firstName, lastName, dob, address, postcode, city, country, state)
			if err != nil {
				return err
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(u)
		},
	}
	cmd.Flags().StringVar(&email, "email", "", "User email (required)")
	cmd.Flags().StringVar(&passwordDigest, "password-digest", "", "User password hash (required)")
	cmd.Flags().StringVar(&referralUID, "referral-uid", "", "Referral UID")
	cmd.Flags().StringVar(&phone, "phone", "", "Phone number")
	cmd.Flags().StringVar(&firstName, "first-name", "", "First name")
	cmd.Flags().StringVar(&lastName, "last-name", "", "Last name")
	cmd.Flags().StringVar(&dob, "dob", "", "Date of birth")
	cmd.Flags().StringVar(&address, "address", "", "Address")
	cmd.Flags().StringVar(&postcode, "postcode", "", "Postcode")
	cmd.Flags().StringVar(&city, "city", "", "City")
	cmd.Flags().StringVar(&country, "country", "", "Country")
	cmd.Flags().StringVar(&state, "state", "", "State")
	_ = cmd.MarkFlagRequired("email")
	_ = cmd.MarkFlagRequired("password-digest")
	return cmd
}
