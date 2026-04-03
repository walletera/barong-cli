package management

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	pkgmgmt "barong-cli/pkg/management"
)

// newManagementClient loads the RSA private key and returns a management API client.
// The key ID and private key file path are resolved from flags first, then env vars.
func newManagementClient(baseURL string) (*pkgmgmt.Client, error) {
	kid := keyID
	if kid == "" {
		kid = os.Getenv("BARONG_MANAGEMENT_KEY_ID")
	}
	if kid == "" {
		return nil, fmt.Errorf("management API key ID is required: set --key-id or BARONG_MANAGEMENT_KEY_ID")
	}

	keyFile := privateKeyFile
	if keyFile == "" {
		keyFile = os.Getenv("BARONG_MANAGEMENT_PRIVATE_KEY_FILE")
	}
	if keyFile == "" {
		return nil, fmt.Errorf("private key file is required: set --private-key-file or BARONG_MANAGEMENT_PRIVATE_KEY_FILE")
	}

	privKey, err := loadPrivateKey(keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load private key: %w", err)
	}

	return pkgmgmt.NewClient(baseURL, kid, privKey), nil
}

func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("no PEM block found in %s", path)
	}
	switch block.Type {
	case "RSA PRIVATE KEY":
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	case "PRIVATE KEY":
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		rsaKey, ok := key.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("PKCS8 key is not RSA")
		}
		return rsaKey, nil
	default:
		return nil, fmt.Errorf("unsupported PEM block type: %s", block.Type)
	}
}
