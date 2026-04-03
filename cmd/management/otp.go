package management

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newOTPCmd(getBaseURL func() string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "otp",
		Short: "OTP operations",
	}
	cmd.AddCommand(newOTPSignCmd(getBaseURL))
	return cmd
}

func newOTPSignCmd(getBaseURL func() string) *cobra.Command {
	var userUID, otpCode string

	cmd := &cobra.Command{
		Use:   "sign",
		Short: "Sign a request with Barong OTP signature",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newManagementClient(getBaseURL())
			if err != nil {
				return err
			}
			if err := client.SignOTP(userUID, otpCode); err != nil {
				return err
			}
			fmt.Fprintln(os.Stderr, "OTP sign successful")
			return nil
		},
	}
	cmd.Flags().StringVar(&userUID, "user-uid", "", "Account UID (required)")
	cmd.Flags().StringVar(&otpCode, "otp-code", "", "Code from Google Authenticator (required)")
	_ = cmd.MarkFlagRequired("user-uid")
	_ = cmd.MarkFlagRequired("otp-code")
	return cmd
}
