package model

import (
	"log"
	"time"
)

type Tx struct {
	ID int64 `json:"id"`
	// CategoryID int64     `json:"categoryId"`
	Category  string    `json:"category"`
	Date      time.Time `json:"date"`
	UserID    int64     `json:"user_id"`
	Amount    int64     `json:"amount"`
	Type      string    `json:"type"`
	Desc      string    `json:"desc"`
	CreatedAt time.Time `json:"createt_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TxReq struct {
	ID int64 `json:"id"`
	// CategoryID int64     `json:"categoryId"`
	Category string `json:"category"`
	Date     string `json:"date"`
	UserID   int64  `json:"user_id"`
	Amount   int64  `json:"amount"`
	Type     string `json:"type"`
	Desc     string `json:"desc"`
}

type TxChart struct {
	Category string `json:"category"`
	Total    int64  `json:"total"`
}

func TxReqToModel(tx TxReq) Tx {
	// const layout = time.RFC3339
	date, err := time.Parse(time.RFC3339, tx.Date)
	if err != nil {
		log.Println(err.Error())
	}
	log.Println(date)
	return Tx{
		ID:       tx.ID,
		Category: tx.Category,
		Date:     date,
		UserID:   tx.UserID,
		Amount:   tx.Amount,
		Type:     tx.Type,
		Desc:     tx.Desc,
	}
}
