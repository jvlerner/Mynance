package auth

import (
	"encoding/json"
	"errors"
	"net/http"
)

type Claims struct {
	UserID    int    `json:"userId"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	ExpiresAt int64  `json:"expiresAt"`
}

var Auth = NewAuthClient()

func ValidateUserToken(userToken string) (*Claims, error) {
	req, err := http.NewRequest(http.MethodGet, Auth.AuthURL+"/auth/service/validate-token", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+Auth.GetToken())
	req.Header.Set("X-User-Token", userToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("unauthorized or invalid user token")
	}

	var result struct {
		Valid     bool   `json:"valid"`
		Error     string `json:"error"`
		UserID    int    `json:"userId"`
		Email     string `json:"email"`
		Role      string `json:"role"`
		ExpiresAt int64  `json:"expiresAt"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if !result.Valid {
		return nil, errors.New("token is not valid: " + result.Error)
	}

	return &Claims{
		UserID:    result.UserID,
		Email:     result.Email,
		Role:      result.Role,
		ExpiresAt: result.ExpiresAt,
	}, nil
}
