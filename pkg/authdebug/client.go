package authdebug

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	auth       authenticator
}

type authenticator interface {
	apply(req *http.Request)
}

// cookieAuth authenticates using a session cookie.
type cookieAuth struct {
	cookies []*http.Cookie
}

func (a *cookieAuth) apply(req *http.Request) {
	for _, c := range a.cookies {
		req.AddCookie(c)
	}
}

// apiKeyAuth authenticates using Barong API key headers.
// Signature: HMAC-SHA256(secret, nonce+kid), hex-encoded.
type apiKeyAuth struct {
	kid    string
	secret string
}

func (a *apiKeyAuth) apply(req *http.Request) {
	nonce := strconv.FormatInt(time.Now().UnixMilli(), 10)
	mac := hmac.New(sha256.New, []byte(a.secret))
	mac.Write([]byte(nonce + a.kid))
	signature := hex.EncodeToString(mac.Sum(nil))

	req.Header.Set("X-Auth-Apikey", a.kid)
	req.Header.Set("X-Auth-Nonce", nonce)
	req.Header.Set("X-Auth-Signature", signature)
}

func NewClientWithCookies(baseURL string, cookies []*http.Cookie) *Client {
	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{},
		auth:       &cookieAuth{cookies: cookies},
	}
}

func NewClientWithAPIKey(baseURL, kid, secret string) *Client {
	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{},
		auth:       &apiKeyAuth{kid: kid, secret: secret},
	}
}

// Debug sends GET /api/v2/auth/{testPath} and returns the HTTP status code and
// response headers. testPath is the path to test, e.g. "api/v1/payments".
func (c *Client) Debug(testPath string) (int, http.Header, error) {
	url := c.baseURL + "/api/v2/auth/" + strings.TrimLeft(testPath, "/")
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, nil, fmt.Errorf("building request: %w", err)
	}
	c.auth.apply(req)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()
	return resp.StatusCode, resp.Header, nil
}
