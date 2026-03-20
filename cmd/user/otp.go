package user

import (
    "encoding/base64"
    "fmt"
    "os"
    "os/exec"
    "runtime"

    "github.com/spf13/cobra"
)

func newOTPCmd(getBaseURL func() string) *cobra.Command {
    cmd := &cobra.Command{
        Use:   "otp",
        Short: "Manage two-factor authentication (2FA)",
    }

    cmd.AddCommand(newOTPGenerateQRCodeCmd(getBaseURL))
    cmd.AddCommand(newOTPEnableCmd(getBaseURL))

    return cmd
}

func newOTPGenerateQRCodeCmd(getBaseURL func() string) *cobra.Command {
    var showSecret bool

    cmd := &cobra.Command{
        Use:   "generate-qrcode",
        Short: "Generate a QR code for setting up 2FA",
        RunE: func(cmd *cobra.Command, args []string) error {
            client, err := newAuthenticatedClient(getBaseURL())
            if err != nil {
                return err
            }
            otp, err := client.GenerateOTPQRCode()
            if err != nil {
                return err
            }

            if showSecret {
                f, err := os.OpenFile(os.TempDir()+"/barong-otp-secret.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
                if err != nil {
                    return fmt.Errorf("failed to create secret file: %w", err)
                }
                secretPath := f.Name()
                if _, err := fmt.Fprintln(f, otp.URL); err != nil {
                    f.Close()
                    return fmt.Errorf("failed to write secret file: %w", err)
                }
                f.Close()
                fmt.Fprintf(os.Stderr, "Secret written to: %s\n", secretPath)
                fmt.Fprintf(os.Stderr, "WARNING: delete this file after saving the secret in your authenticator app.\n")
                return openFile(secretPath)
            }

            imgData, err := base64.StdEncoding.DecodeString(otp.Barcode)
            if err != nil {
                return fmt.Errorf("failed to decode QR code image: %w", err)
            }
            f, err := os.OpenFile(os.TempDir()+"/barong-otp-qrcode.png", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
            if err != nil {
                return fmt.Errorf("failed to create QR code file: %w", err)
            }
            imgPath := f.Name()
            if _, err := f.Write(imgData); err != nil {
                f.Close()
                return fmt.Errorf("failed to write QR code image: %w", err)
            }
            f.Close()
            fmt.Fprintf(os.Stderr, "QR code saved to: %s\n", imgPath)
            fmt.Fprintf(os.Stderr, "WARNING: delete this file after scanning the QR code.\n")
            fmt.Fprintf(os.Stderr, "Tip: run with --show-secret to get the secret key as text instead.\n")
            return openFile(imgPath)
        },
    }

    cmd.Flags().BoolVar(&showSecret, "show-secret", false, "Write the secret key to a file instead of showing the QR code")
    return cmd
}

func openFile(path string) error {
    var opener string
    switch runtime.GOOS {
    case "darwin":
        opener = "open"
    case "windows":
        opener = "start"
    default:
        opener = "xdg-open"
    }
    return exec.Command(opener, path).Start()
}

func newOTPEnableCmd(getBaseURL func() string) *cobra.Command {
    var code string

    cmd := &cobra.Command{
        Use:   "enable",
        Short: "Enable 2FA using a code from your authenticator app",
        RunE: func(cmd *cobra.Command, args []string) error {
            client, err := newAuthenticatedClient(getBaseURL())
            if err != nil {
                return err
            }
            if err := client.EnableOTP(code); err != nil {
                return err
            }
            fmt.Println("2FA enabled")
            return nil
        },
    }

    cmd.Flags().StringVar(&code, "code", "", "Code from Google Authenticator (required)")
    _ = cmd.MarkFlagRequired("code")

    return cmd
}
