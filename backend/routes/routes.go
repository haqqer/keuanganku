package routes

import (
	"net/http"

	"github.com/haqqer/keuanganku/handler"
	"github.com/haqqer/keuanganku/middleware"
)

func Router() *http.ServeMux {
	router := http.NewServeMux()

	authHandler := handler.AuthHandler{}
	// userHandler := handler.UsersHandler{}
	txHandler := handler.TxHandler{}

	router.HandleFunc("GET /", handler.HomeHandler)

	router.HandleFunc("GET /auth/google", authHandler.Login)
	router.HandleFunc("GET /callback", authHandler.Callback)
	router.HandleFunc("POST /refresh-token", authHandler.RefreshToken)
	router.HandleFunc("GET /auth/me", middleware.CheckToken(authHandler.AuthMe))

	router.HandleFunc("GET /api/chart", middleware.CheckToken(txHandler.GetChart))
	router.HandleFunc("GET /api/tx", middleware.CheckToken(txHandler.Get))
	router.HandleFunc("POST /api/tx", middleware.CheckToken(txHandler.Create))
	router.HandleFunc("DELETE /api/tx/{id}", middleware.CheckToken(txHandler.DeleteByID))
	router.HandleFunc("PUT /api/tx/{id}", middleware.CheckToken(txHandler.UpdateByID))

	return router
}
