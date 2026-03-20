package user

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	cookies    []*http.Cookie // non-nil only for authenticated clients
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{},
	}
}

func NewAuthenticatedClient(baseURL string, cookies []*http.Cookie) *Client {
	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{},
		cookies:    cookies,
	}
}

func (c *Client) CreateUser(email, password, username, refid string) (*UserWithFullInfo, error) {
	form := url.Values{
		"email":    {email},
		"password": {password},
	}
	if username != "" {
		form.Set("username", username)
	}
	if refid != "" {
		form.Set("refid", refid)
	}
	resp, err := c.post("/api/v1/auth/identity/users", form, false)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return decodeUserWithFullInfo(resp)
}

func (c *Client) Login(email, password, otpCode string) (*UserWithFullInfo, []*http.Cookie, error) {
	form := url.Values{
		"email":    {email},
		"password": {password},
	}
	if otpCode != "" {
		form.Set("otp_code", otpCode)
	}
	req, err := http.NewRequest(http.MethodPost, c.baseURL+"/api/v1/auth/identity/sessions", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	u, err := decodeUserWithFullInfo(resp)
	if err != nil {
		return nil, nil, err
	}
	return u, resp.Cookies(), nil
}

func (c *Client) Logout(cookies []*http.Cookie) error {
	req, err := http.NewRequest(http.MethodDelete, c.baseURL+"/api/v1/auth/identity/sessions", nil)
	if err != nil {
		return err
	}
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("logout failed (%d): %s", resp.StatusCode, body)
	}
	return nil
}

func (c *Client) GetMe() (*UserWithFullInfo, error) {
	resp, err := c.get("/api/v1/auth/resource/users/me", nil, true)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return decodeUserWithFullInfo(resp)
}

func (c *Client) GenerateOTPQRCode() (*OTPQRCode, error) {
	resp, err := c.post("/api/v1/auth/resource/otp/generate_qrcode", url.Values{}, true)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, body)
	}
	var result struct {
		Data OTPQRCode `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &result.Data, nil
}

func (c *Client) EnableOTP(code string) error {
	resp, err := c.post("/api/v1/auth/resource/otp/enable", url.Values{"code": {code}}, true)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error %d: %s", resp.StatusCode, body)
	}
	return nil
}

// --- helpers ---

func (c *Client) post(path string, form url.Values, authenticated bool) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, c.baseURL+path, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if authenticated {
		for _, cookie := range c.cookies {
			req.AddCookie(cookie)
		}
	}
	return c.httpClient.Do(req)
}

func (c *Client) get(path string, queryParams url.Values, authenticated bool) (*http.Response, error) {
	fullURL := c.baseURL + path
	if len(queryParams) > 0 {
		fullURL += "?" + queryParams.Encode()
	}
	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, err
	}
	if authenticated {
		for _, cookie := range c.cookies {
			req.AddCookie(cookie)
		}
	}
	return c.httpClient.Do(req)
}

func decodeUserWithFullInfo(resp *http.Response) (*UserWithFullInfo, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, body)
	}
	var u UserWithFullInfo
	if err := json.Unmarshal(body, &u); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &u, nil
}
