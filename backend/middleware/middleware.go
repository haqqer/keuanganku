package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/haqqer/keuanganku/repo"
	"github.com/haqqer/keuanganku/utils/auth"
	"github.com/haqqer/keuanganku/utils/response"
	"golang.org/x/oauth2"
)

type Middleware func(http.Handler) http.Handler

const AuthUser = "middleware.auth.user"

func Loader(md ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(md) - 1; i >= 0; i-- {
			x := md[i]
			next = x(next)
		}
		return next
	}
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Println(r.Method, r.URL.Path, time.Since(start))
	})
}

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		if r.Method == "OPTIONS" {
			w.Write([]byte("allowed"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func CheckToken(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var token = new(oauth2.Token)
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

		token.AccessToken = srcToken[1]
		token.Expiry = time.Time{}
		authUser, err := auth.GetUserInfo(r.Context(), token)
		if err != nil {
			response.Error(w, http.StatusUnauthorized, err.Error())
			return
		}
		userRepo := repo.UserRepo{}
		user, err := userRepo.GetByGoogleID(r.Context(), authUser.ID)
		if err != nil {
			response.Error(w, http.StatusUnauthorized, "user not registered yet!")
			return
		}
		if user.GoogleID == "" {
			response.Error(w, http.StatusUnauthorized, "user not registered yet!")
			return
		}
		ctx := context.WithValue(r.Context(), AuthUser, user)
		req := r.WithContext(ctx)
		f(w, req)
	}
}
