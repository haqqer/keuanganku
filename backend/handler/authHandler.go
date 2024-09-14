package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/haqqer/keuanganku/middleware"
	"github.com/haqqer/keuanganku/model"
	"github.com/haqqer/keuanganku/repo"
	"github.com/haqqer/keuanganku/utils/auth"
	"github.com/haqqer/keuanganku/utils/response"
	"golang.org/x/oauth2"
)

type AuthHandler struct{}

func (h *AuthHandler) AuthMe(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.AuthUser).(model.User)
	response.Success(w, user)

}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	authCode := os.Getenv("AUTH_CODE")
	url := auth.AuthConfig.AuthCodeURL(authCode, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	http.Redirect(w, r, url, http.StatusSeeOther)
	response.Success(w, url)
}

func (h *AuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	userRepo := repo.UserRepo{}
	authCode := os.Getenv("AUTH_CODE")
	state := r.URL.Query().Get("state")
	if state != authCode {
		response.Error(w, http.StatusBadRequest, fmt.Sprint(w, "code state not same"))
		return
	}
	code := r.URL.Query().Get("code")
	authGoogle := auth.GoogleConfig()

	token, err := authGoogle.Exchange(r.Context(), code)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, fmt.Sprint("Code-Token Exchange Failed!", err.Error()))
		return
	}
	authUser, err := auth.GetUserInfo(r.Context(), token)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, fmt.Sprint("error get user info : ", err.Error()))
		return
	}

	user := model.User{
		Email:        authUser.Email,
		Username:     "-",
		Name:         authUser.Name,
		GoogleID:     authUser.ID,
		PictureUrl:   authUser.Picture,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
	}
	existingUser, err := userRepo.GetByGoogleID(r.Context(), user.GoogleID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if existingUser.GoogleID == "" {
		if err := userRepo.Create(r.Context(), user); err != nil {
			response.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	url := os.Getenv("CLIENT_URL")
	clientUrl := fmt.Sprintf("%s/?access_token=%s&refresh_token=%s", url, token.AccessToken, token.RefreshToken)
	http.Redirect(w, r, clientUrl, http.StatusSeeOther)
	// response.Success(w, map[string]interface{}{
	// 	"user_data": authUser,
	// 	"token":     token,
	// })
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var token = new(oauth2.Token)
	var bodyToken *oauth2.Token
	var authUserData *auth.AuthUserData
	err := json.NewDecoder(r.Body).Decode(&bodyToken)
	if err != nil {
		fmt.Fprintf(w, "Failed to parse")
		return
	}
	accessToken := r.Header.Get("authorization")
	if accessToken == "" {
		response.Error(w, http.StatusUnauthorized, "token empty!")
		return
	}
	srcToken := strings.Split(accessToken, " ")
	if len(srcToken) < 1 {
		response.Error(w, http.StatusUnauthorized, "token invalid!")
		return
	}

	token.RefreshToken = bodyToken.RefreshToken
	newToken, err := auth.RefreshToken(r.Context(), token)
	if err != nil {
		http.Error(w, fmt.Sprint("error obtain new token : ", err.Error()), http.StatusInternalServerError)
		return
	}

	authUserData, err = auth.GetUserInfo(r.Context(), newToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, map[string]interface{}{
		"user_data": authUserData,
		"token":     newToken,
	})
}
