package management

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
)

func newPhonesCmd(getBaseURL func() string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "phones",
		Short: "Phone number management operations",
	}
	cmd.AddCommand(newPhoneCreateCmd(getBaseURL))
	cmd.AddCommand(newPhoneGetCmd(getBaseURL))
	cmd.AddCommand(newPhoneDeleteCmd(getBaseURL))
	return cmd
}

func newPhoneCreateCmd(getBaseURL func() string) *cobra.Command {
	var uid, number string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a phone number for a user",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newManagementClient(getBaseURL())
			if err != nil {
				return err
			}
			phone, err := client.CreatePhone(uid, number)
			if err != nil {
				return err
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(phone)
		},
	}
	cmd.Flags().StringVar(&uid, "uid", "", "User UID (required)")
	cmd.Flags().StringVar(&number, "number", "", "Phone number (required)")
	_ = cmd.MarkFlagRequired("uid")
	_ = cmd.MarkFlagRequired("number")
	return cmd
}

func newPhoneGetCmd(getBaseURL func() string) *cobra.Command {
	var uid string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get phone numbers for a user",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newManagementClient(getBaseURL())
			if err != nil {
				return err
			}
			phones, err := client.GetPhones(uid)
			if err != nil {
				return err
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(phones)
		},
	}
	cmd.Flags().StringVar(&uid, "uid", "", "User UID (required)")
	_ = cmd.MarkFlagRequired("uid")
	return cmd
}

func newPhoneDeleteCmd(getBaseURL func() string) *cobra.Command {
	var uid, number string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a phone number for a user",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newManagementClient(getBaseURL())
			if err != nil {
				return err
			}
			phone, err := client.DeletePhone(uid, number)
			if err != nil {
				return err
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(phone)
		},
	}
	cmd.Flags().StringVar(&uid, "uid", "", "User UID (required)")
	cmd.Flags().StringVar(&number, "number", "", "Phone number (required)")
	_ = cmd.MarkFlagRequired("uid")
	_ = cmd.MarkFlagRequired("number")
	return cmd
}
