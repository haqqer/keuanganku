package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/haqqer/keuanganku/middleware"
	"github.com/haqqer/keuanganku/model"
	"github.com/haqqer/keuanganku/repo"
	"github.com/haqqer/keuanganku/utils/response"
)

type TxHandler struct{}

func (h *TxHandler) GetChart(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.AuthUser).(model.User)
	month := r.URL.Query().Get("month")
	n, err := strconv.Atoi(month)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	txType := r.URL.Query().Get("type")
	txRepo := repo.TxRepo{}
	txs, err := txRepo.GetChartByUser(r.Context(), user.ID, n, txType)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, fmt.Sprintf("error data fetch : %s", err.Error()))
		return
	}

	response.Success(w, txs)
}
func (h *TxHandler) Get(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.AuthUser).(model.User)
	month := r.URL.Query().Get("month")
	if month == "" {
		month = "0"
	}
	qMonth, err := strconv.Atoi(month)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	year := r.URL.Query().Get("year")
	if year == "" {
		year = "0"
	}
	qYear, err := strconv.Atoi(year)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	txRepo := repo.TxRepo{}
	// query := repo.NewQuery(10, 0, "date", false)
	// txs, err := txRepo.GetAll(r.Context(), user.ID, query)
	txs, err := txRepo.GetByMonth(r.Context(), user.ID, qMonth, qYear)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, fmt.Sprintf("error data fetch : %s", err.Error()))
		return
	}

	response.Success(w, txs)
}

func (h *TxHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.AuthUser).(model.User)
	txRepo := repo.TxRepo{}

	var tx model.Tx
	err := json.NewDecoder(r.Body).Decode(&tx)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, fmt.Sprintf("error parse data : %s", err.Error()))
		return
	}

	tx.UserID = user.ID

	if err := txRepo.Create(r.Context(), tx); err != nil {
		response.Error(w, http.StatusInternalServerError, fmt.Sprintf("error input data : %s", err.Error()))
		return
	}
	response.Success(w, tx)
}

func (h *TxHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	txID := r.PathValue("id")
	if txID == "" {
		response.Error(w, http.StatusBadRequest, "ID params empty")
		return
	}
	n, _ := strconv.ParseInt(txID, 10, 64)
	txRepo := repo.TxRepo{}

	if err := txRepo.DeleteByID(r.Context(), n); err != nil {
		response.Error(w, http.StatusInternalServerError, fmt.Sprintf("error delete data : %s", err.Error()))
		return
	}
	response.Success(w, map[string]string{})
}

func (h *TxHandler) UpdateByID(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.AuthUser).(model.User)
	txID := r.PathValue("id")
	if txID == "" {
		response.Error(w, http.StatusBadRequest, "ID params empty")
		return
	}
	ID, _ := strconv.ParseInt(txID, 10, 64)
	txRepo := repo.TxRepo{}

	var tx model.Tx
	err := json.NewDecoder(r.Body).Decode(&tx)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, fmt.Sprintf("error parse data : %s", err.Error()))
		return
	}
	tx.UserID = user.ID
	tx.ID = ID

	if err := txRepo.UpdateByID(r.Context(), tx); err != nil {
		response.Error(w, http.StatusInternalServerError, fmt.Sprintf("error input data : %s", err.Error()))
		return
	}
	response.Success(w, tx)
}
