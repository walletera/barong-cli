package session

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Session struct {
	Cookies []cookieInfo `json:"cookies"`
}

type cookieInfo struct {
	Name    string    `json:"name"`
	Value   string    `json:"value"`
	Expires time.Time `json:"expires,omitempty"`
}

func sessionPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".barong-cli", "session.json"), nil
}

func Save(cookies []*http.Cookie) error {
	path, err := sessionPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	s := Session{}
	for _, c := range cookies {
		s.Cookies = append(s.Cookies, cookieInfo{Name: c.Name, Value: c.Value, Expires: c.Expires})
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

func Load() ([]*http.Cookie, error) {
	path, err := sessionPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, errors.New("not logged in: run 'barong-cli user login' first")
	}
	if err != nil {
		return nil, err
	}
	var s Session
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	cookies := make([]*http.Cookie, len(s.Cookies))
	for i, ci := range s.Cookies {
		if !ci.Expires.IsZero() && time.Now().After(ci.Expires) {
			return nil, errors.New("session expired: run 'barong-cli user login' again")
		}
		cookies[i] = &http.Cookie{Name: ci.Name, Value: ci.Value, Expires: ci.Expires}
	}
	return cookies, nil
}

func Delete() error {
	path, err := sessionPath()
	if err != nil {
		return err
	}
	err = os.Remove(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
