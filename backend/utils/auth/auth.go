package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthUserData struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

var AuthConfig oauth2.Config

func GoogleConfig() oauth2.Config {
	cfg := oauth2.Config{
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint: google.Endpoint,
	}

	return cfg
}

func RefreshToken(ctx context.Context, token *oauth2.Token) (*oauth2.Token, error) {
	newToken, err := AuthConfig.TokenSource(ctx, token).Token()
	if err != nil {
		return nil, err
	}
	return newToken, nil
}

func GetUserInfo(ctx context.Context, token *oauth2.Token) (*AuthUserData, error) {
	var authUserData *AuthUserData
	client := AuthConfig.Client(ctx, token)
	log.Println(token.AccessToken)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	// resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("error get userinfo : %s", err.Error())
	}
	if resp.StatusCode == http.StatusUnauthorized {
		p := map[string]interface{}{}
		json.NewDecoder(resp.Body).Decode(&p)
		errorDetail := p["error"].(map[string]interface{})
		return nil, fmt.Errorf("error parsing authUserdata error : %s", errorDetail["message"])
	}
	err = json.NewDecoder(resp.Body).Decode(&authUserData)
	if err != nil {
		return nil, fmt.Errorf("error parsing authUserdata error : %s", err.Error())
	}
	defer resp.Body.Close()
	return authUserData, nil
}
