package management

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
)

func newServiceAccountsCmd(getBaseURL func() string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service-accounts",
		Short: "Service account management operations",
	}
	cmd.AddCommand(newServiceAccountCreateCmd(getBaseURL))
	cmd.AddCommand(newServiceAccountGetCmd(getBaseURL))
	cmd.AddCommand(newServiceAccountListCmd(getBaseURL))
	cmd.AddCommand(newServiceAccountDeleteCmd(getBaseURL))
	return cmd
}

func newServiceAccountCreateCmd(getBaseURL func() string) *cobra.Command {
	var ownerUID, role, serviceAccountUID, serviceAccountEmail string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a service account",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newManagementClient(getBaseURL())
			if err != nil {
				return err
			}
			sa, err := client.CreateServiceAccount(ownerUID, role, serviceAccountUID, serviceAccountEmail)
			if err != nil {
				return err
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(sa)
		},
	}
	cmd.Flags().StringVar(&ownerUID, "owner-uid", "", "Owner UID (required)")
	cmd.Flags().StringVar(&role, "role", "", "Service account role (required)")
	cmd.Flags().StringVar(&serviceAccountUID, "service-account-uid", "", "Service account UID")
	cmd.Flags().StringVar(&serviceAccountEmail, "service-account-email", "", "Service account email")
	_ = cmd.MarkFlagRequired("owner-uid")
	_ = cmd.MarkFlagRequired("role")
	return cmd
}

func newServiceAccountGetCmd(getBaseURL func() string) *cobra.Command {
	var uid, email string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get service account information",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newManagementClient(getBaseURL())
			if err != nil {
				return err
			}
			sa, err := client.GetServiceAccount(uid, email)
			if err != nil {
				return err
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(sa)
		},
	}
	cmd.Flags().StringVar(&uid, "uid", "", "Service account UID")
	cmd.Flags().StringVar(&email, "email", "", "Service account email")
	return cmd
}

func newServiceAccountListCmd(getBaseURL func() string) *cobra.Command {
	var ownerUID, ownerEmail string
	var page, limit int64

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List service accounts",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newManagementClient(getBaseURL())
			if err != nil {
				return err
			}
			accounts, err := client.ListServiceAccounts(ownerUID, ownerEmail, page, limit)
			if err != nil {
				return err
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(accounts)
		},
	}
	cmd.Flags().StringVar(&ownerUID, "owner-uid", "", "Filter by owner UID")
	cmd.Flags().StringVar(&ownerEmail, "owner-email", "", "Filter by owner email")
	cmd.Flags().Int64Var(&page, "page", 0, "Page number (default 1)")
	cmd.Flags().Int64Var(&limit, "limit", 0, "Results per page (max 100)")
	return cmd
}

func newServiceAccountDeleteCmd(getBaseURL func() string) *cobra.Command {
	var uid string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a service account",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newManagementClient(getBaseURL())
			if err != nil {
				return err
			}
			sa, err := client.DeleteServiceAccount(uid)
			if err != nil {
				return err
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(sa)
		},
	}
	cmd.Flags().StringVar(&uid, "uid", "", "Service account UID (required)")
	_ = cmd.MarkFlagRequired("uid")
	return cmd
}
