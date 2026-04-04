package user

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newAPIKeyCmd(getBaseURL func() string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api-key",
		Short: "Manage API keys for the current account",
	}

	cmd.AddCommand(newAPIKeyListCmd(getBaseURL))
	cmd.AddCommand(newAPIKeyCreateCmd(getBaseURL))
	cmd.AddCommand(newAPIKeyUpdateCmd(getBaseURL))
	cmd.AddCommand(newAPIKeyDeleteCmd(getBaseURL))

	return cmd
}

func newAPIKeyListCmd(getBaseURL func() string) *cobra.Command {
	var page, limit int
	var orderBy, ordering, serviceAccountUID string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all API keys for the current account",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newAuthenticatedClient(getBaseURL())
			if err != nil {
				return err
			}
			keys, err := client.ListAPIKeys(page, limit, orderBy, ordering, serviceAccountUID)
			if err != nil {
				return err
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(keys)
		},
	}

	cmd.Flags().IntVar(&page, "page", 0, "Page number")
	cmd.Flags().IntVar(&limit, "limit", 0, "Results per page (max 100)")
	cmd.Flags().StringVar(&orderBy, "order-by", "", "Field to sort by")
	cmd.Flags().StringVar(&ordering, "ordering", "", "Sort order (asc or desc)")
	cmd.Flags().StringVar(&serviceAccountUID, "service-account-uid", "", "Service account UID (lists keys for that service account)")

	return cmd
}

func newAPIKeyCreateCmd(getBaseURL func() string) *cobra.Command {
	var algorithm, scope, totpCode, serviceAccountUID string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new API key",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newAuthenticatedClient(getBaseURL())
			if err != nil {
				return err
			}
			key, err := client.CreateAPIKey(algorithm, scope, totpCode, serviceAccountUID)
			if err != nil {
				return err
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(key)
		},
	}

	cmd.Flags().StringVar(&algorithm, "algorithm", "", "API key algorithm (required)")
	cmd.Flags().StringVar(&scope, "scope", "", "Comma-separated scopes")
	cmd.Flags().StringVar(&totpCode, "totp-code", "", "Code from Google Authenticator (required)")
	cmd.Flags().StringVar(&serviceAccountUID, "service-account-uid", "", "Service account UID (creates key for that service account)")
	_ = cmd.MarkFlagRequired("algorithm")
	_ = cmd.MarkFlagRequired("totp-code")

	return cmd
}

func newAPIKeyUpdateCmd(getBaseURL func() string) *cobra.Command {
	var scope, state, totpCode, serviceAccountUID string

	cmd := &cobra.Command{
		Use:   "update <kid>",
		Short: "Update an API key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newAuthenticatedClient(getBaseURL())
			if err != nil {
				return err
			}
			key, err := client.UpdateAPIKey(args[0], scope, state, totpCode, serviceAccountUID)
			if err != nil {
				return err
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(key)
		},
	}

	cmd.Flags().StringVar(&scope, "scope", "", "Comma-separated scopes")
	cmd.Flags().StringVar(&state, "state", "", "Key state (active or disabled)")
	cmd.Flags().StringVar(&totpCode, "totp-code", "", "Code from Google Authenticator (required)")
	cmd.Flags().StringVar(&serviceAccountUID, "service-account-uid", "", "Service account UID (updates key for that service account)")
	_ = cmd.MarkFlagRequired("totp-code")

	return cmd
}

func newAPIKeyDeleteCmd(getBaseURL func() string) *cobra.Command {
	var totpCode, serviceAccountUID string

	cmd := &cobra.Command{
		Use:   "delete <kid>",
		Short: "Delete an API key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newAuthenticatedClient(getBaseURL())
			if err != nil {
				return err
			}
			if err := client.DeleteAPIKey(args[0], totpCode, serviceAccountUID); err != nil {
				return err
			}
			fmt.Fprintln(os.Stderr, "API key deleted")
			return nil
		},
	}

	cmd.Flags().StringVar(&totpCode, "totp-code", "", "Code from Google Authenticator (required)")
	cmd.Flags().StringVar(&serviceAccountUID, "service-account-uid", "", "Service account UID (deletes key for that service account)")
	_ = cmd.MarkFlagRequired("totp-code")

	return cmd
}
