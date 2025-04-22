package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

type AuthClient struct {
	Token           string
	ExpiresAt       int64
	mutex           sync.Mutex
	ServiceEmail    string
	ServicePassword string
	AuthURL         string
}

func NewAuthClient() *AuthClient {
	return &AuthClient{
		ServiceEmail:    os.Getenv("SERVICE_EMAIL"),
		ServicePassword: os.Getenv("SERVICE_PASSWORD"),
		AuthURL:         os.Getenv("AUTH_URL"),
	}
}

func (a *AuthClient) Login() error {
	body := map[string]string{
		"email":    a.ServiceEmail,
		"password": a.ServicePassword,
	}
	b, _ := json.Marshal(body)

	resp, err := http.Post(a.AuthURL+"/auth/service/login", "application/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result struct {
		Token     string `json:"token"`
		ExpiresAt int64  `json:"expiresAt"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	a.mutex.Lock()
	a.Token = result.Token
	a.ExpiresAt = result.ExpiresAt
	a.mutex.Unlock()

	return nil
}

func (a *AuthClient) GetToken() string {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if a.Token == "" {
		if err := a.Login(); err != nil {
			fmt.Println("[ERROR] [AUTH] Error logging in:", err)
			return ""
		}
	}
	if time.Until(time.Unix(a.ExpiresAt, 0)) < 2*time.Minute {
		go a.Login()
	}

	a.mutex.Lock()
	defer a.mutex.Unlock()
	return a.Token
}
