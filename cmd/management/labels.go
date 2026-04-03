package management

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
)

func newLabelsCmd(getBaseURL func() string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "labels",
		Short: "Label management operations",
	}
	cmd.AddCommand(newLabelCreateCmd(getBaseURL))
	cmd.AddCommand(newLabelUpdateCmd(getBaseURL))
	cmd.AddCommand(newLabelDeleteCmd(getBaseURL))
	cmd.AddCommand(newLabelListCmd(getBaseURL))
	cmd.AddCommand(newLabelFilterUsersCmd(getBaseURL))
	return cmd
}

func newLabelCreateCmd(getBaseURL func() string) *cobra.Command {
	var userUID, key, value, description string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a label with 'private' scope for a user",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newManagementClient(getBaseURL())
			if err != nil {
				return err
			}
			label, err := client.CreateLabel(userUID, key, value, description)
			if err != nil {
				return err
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(label)
		},
	}
	cmd.Flags().StringVar(&userUID, "user-uid", "", "User UID (required)")
	cmd.Flags().StringVar(&key, "key", "", "Label key (required)")
	cmd.Flags().StringVar(&value, "value", "", "Label value (required)")
	cmd.Flags().StringVar(&description, "description", "", "Label description")
	_ = cmd.MarkFlagRequired("user-uid")
	_ = cmd.MarkFlagRequired("key")
	_ = cmd.MarkFlagRequired("value")
	return cmd
}

func newLabelUpdateCmd(getBaseURL func() string) *cobra.Command {
	var userUID, key, value, description string
	var replace bool

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a label with 'private' scope",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newManagementClient(getBaseURL())
			if err != nil {
				return err
			}
			label, err := client.UpdateLabel(userUID, key, value, description, replace)
			if err != nil {
				return err
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(label)
		},
	}
	cmd.Flags().StringVar(&userUID, "user-uid", "", "User UID (required)")
	cmd.Flags().StringVar(&key, "key", "", "Label key (required)")
	cmd.Flags().StringVar(&value, "value", "", "Label value (required)")
	cmd.Flags().StringVar(&description, "description", "", "Label description")
	cmd.Flags().BoolVar(&replace, "replace", false, "Create label if it does not exist")
	_ = cmd.MarkFlagRequired("user-uid")
	_ = cmd.MarkFlagRequired("key")
	_ = cmd.MarkFlagRequired("value")
	return cmd
}

func newLabelDeleteCmd(getBaseURL func() string) *cobra.Command {
	var userUID, key string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a label with 'private' scope",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newManagementClient(getBaseURL())
			if err != nil {
				return err
			}
			if err := client.DeleteLabel(userUID, key); err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&userUID, "user-uid", "", "User UID (required)")
	cmd.Flags().StringVar(&key, "key", "", "Label key (required)")
	_ = cmd.MarkFlagRequired("user-uid")
	_ = cmd.MarkFlagRequired("key")
	return cmd
}

func newLabelListCmd(getBaseURL func() string) *cobra.Command {
	var userUID string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List labels for a user",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newManagementClient(getBaseURL())
			if err != nil {
				return err
			}
			labels, err := client.ListLabels(userUID)
			if err != nil {
				return err
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(labels)
		},
	}
	cmd.Flags().StringVar(&userUID, "user-uid", "", "User UID (required)")
	_ = cmd.MarkFlagRequired("user-uid")
	return cmd
}

func newLabelFilterUsersCmd(getBaseURL func() string) *cobra.Command {
	var key, value, scope string
	var extended bool
	var page, limit int64

	cmd := &cobra.Command{
		Use:   "filter-users",
		Short: "Get users filtered by label attributes",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newManagementClient(getBaseURL())
			if err != nil {
				return err
			}
			users, err := client.FilterUsersByLabel(key, value, scope, extended, page, limit)
			if err != nil {
				return err
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(users)
		},
	}
	cmd.Flags().StringVar(&key, "key", "", "Label key (required)")
	cmd.Flags().StringVar(&value, "value", "", "Label value")
	cmd.Flags().StringVar(&scope, "scope", "", "Label scope")
	cmd.Flags().BoolVar(&extended, "extended", false, "Return full user information")
	cmd.Flags().Int64Var(&page, "page", 0, "Page number (default 1)")
	cmd.Flags().Int64Var(&limit, "limit", 0, "Users per page (max 100)")
	_ = cmd.MarkFlagRequired("key")
	return cmd
}
