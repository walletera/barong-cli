package user

import (
	"barong-cli/internal/session"
	pkguser "barong-cli/pkg/user"
)

// newAuthenticatedClient loads the saved session cookies and returns a client
// ready to call /api/v1/auth/resource endpoints using the session cookie directly.
func newAuthenticatedClient(baseURL string) (*pkguser.Client, error) {
	cookies, err := session.Load()
	if err != nil {
		return nil, err
	}
	return pkguser.NewAuthenticatedClient(baseURL, cookies), nil
}
