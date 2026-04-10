package authdebug

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"

	"barong-cli/internal/session"
	"barong-cli/pkg/authdebug"

	"github.com/spf13/cobra"
)

func NewAuthDebugCmd(getBaseURL func() string) *cobra.Command {
	var apiKeyKID, apiKeySecret string

	cmd := &cobra.Command{
		Use:   "auth-debug <test-path>",
		Short: "Test Barong /api/v2/auth/{path} and print response headers",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			testPath := args[0]
			client, err := newClient(getBaseURL(), apiKeyKID, apiKeySecret)
			if err != nil {
				return err
			}
			status, headers, err := client.Debug(testPath)
			if err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "Status: %d %s\n", status, http.StatusText(status))
			printHeaders(headers)
			if bearer := headers.Get("Authorization"); bearer != "" {
				printJWTPayload(bearer)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&apiKeyKID, "api-key-kid", "", "API key ID (uses api-key auth instead of session cookie)")
	cmd.Flags().StringVar(&apiKeySecret, "api-key-secret", "", "API key secret (required when --api-key-kid is set)")

	return cmd
}

func newClient(baseURL, kid, secret string) (*authdebug.Client, error) {
	if kid != "" {
		if secret == "" {
			return nil, fmt.Errorf("--api-key-secret is required when --api-key-kid is set")
		}
		return authdebug.NewClientWithAPIKey(baseURL, kid, secret), nil
	}
	cookies, err := session.Load()
	if err != nil {
		return nil, fmt.Errorf("no active session (login first or provide --api-key-kid/--api-key-secret): %w", err)
	}
	return authdebug.NewClientWithCookies(baseURL, cookies), nil
}

func printJWTPayload(authHeader string) {
	token := strings.TrimPrefix(authHeader, "Bearer ")
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		fmt.Fprintf(os.Stderr, "Authorization header is not a valid JWT\n")
		return
	}
	decoded, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to decode JWT payload: %v\n", err)
		return
	}
	var buf bytes.Buffer
	if err := json.Indent(&buf, decoded, "", "  "); err != nil {
		fmt.Fprintf(os.Stderr, "failed to pretty-print JWT payload: %v\n", err)
		return
	}
	fmt.Printf("\nAuthorization JWT payload:\n%s\n", buf.String())
}

func printHeaders(headers http.Header) {
	names := make([]string, 0, len(headers))
	for name := range headers {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		for _, val := range headers[name] {
			fmt.Printf("%s: %s\n", name, val)
		}
	}
}
