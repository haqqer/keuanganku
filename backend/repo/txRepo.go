package repo

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/haqqer/keuanganku/database"
	"github.com/haqqer/keuanganku/model"
	"github.com/jackc/pgx/v5"
)

type TxRepo struct{}

func (r *TxRepo) Create(ctx context.Context, tx model.Tx) error {
	db := database.DB

	now := time.Now()
	query := `INSERT INTO public."tx" ("category","date","user_id", "amount", "type", "desc", "created_at", "updated_at") VALUES (@category, @date, @user_id, @amount, @type, @desc, @created_at, @updated_at)`
	args := pgx.NamedArgs{
		"category":   tx.Category,
		"date":       tx.Date,
		"user_id":    tx.UserID,
		"amount":     tx.Amount,
		"type":       tx.Type,
		"desc":       tx.Desc,
		"created_at": now,
		"updated_at": now,
	}
	_, err := db.Exec(ctx, query, args)
	if err != nil {
		return err
	}
	return nil
}

func (r *TxRepo) GetAll(ctx context.Context, userId int64, query Query) ([]model.Tx, error) {
	db := database.DB
	var order = "ASC"
	if !query.Asc {
		order = "DESC"
	}

	q := fmt.Sprintf(`SELECT * FROM public."tx" WHERE "user_id" = @user_id ORDER BY "%s" %s LIMIT %d OFFSET %d`, query.OrderBy, order, query.Limit, query.Offset)
	args := pgx.NamedArgs{
		"user_id": userId,
	}
	rows, err := db.Query(ctx, q, args)
	if err != nil {
		if err == pgx.ErrNoRows {
			return []model.Tx{}, nil
		}
		return nil, err
	}
	txs, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Tx])
	if err != nil {
		return nil, err
	}

	return txs, nil
}

func (r *TxRepo) GetByMonth(ctx context.Context, userId int64, month, year int) ([]model.Tx, error) {
	db := database.DB
	t := time.Now()
	qMonth := fmt.Sprintf(`AND EXTRACT(MONTH FROM "date") = %d`, t.Month())
	qYear := fmt.Sprintf(`AND EXTRACT(YEAR FROM "date") = %d`, t.Year())

	if month > 0 {
		qMonth = fmt.Sprintf(`AND EXTRACT(MONTH FROM "date") = %d`, month)
	}
	if year > 0 {
		qYear = fmt.Sprintf(`AND EXTRACT(YEAR FROM "date") = %d`, year)
	}

	q := `SELECT * FROM public."tx" WHERE "user_id" = @user_id ` + qMonth + ` ` + qYear + ` ORDER BY "date" DESC`
	log.Println(q)
	args := pgx.NamedArgs{
		"user_id": userId,
	}
	rows, err := db.Query(ctx, q, args)
	if err != nil {
		log.Println("error disini")
		if err == pgx.ErrNoRows {
			return []model.Tx{}, nil
		}
		return nil, err
	}
	txs, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Tx])
	if err != nil {
		log.Println("error disini")
		return nil, err
	}

	return txs, nil
}

func (r *TxRepo) GetByID(ctx context.Context, userId int64) (model.Tx, error) {
	db := database.DB
	var tx model.Tx

	query := `SELECT * FROM public."tx" WHERE "user_id" = @user_id`
	args := pgx.NamedArgs{
		"user_id": userId,
	}
	err := db.QueryRow(ctx, query, args).Scan(&tx.ID, &tx.Category, &tx.UserID, &tx.Amount, &tx.Type, &tx.Desc, &tx.CreatedAt, &tx.UpdatedAt)
	if err != nil {
		return tx, err
	}
	return tx, nil
}

func (r *TxRepo) GetChartByUser(ctx context.Context, userId int64, month int, txType string) ([]model.TxChart, error) {
	db := database.DB
	var txs []model.TxChart
	t := time.Now()

	qType := ""
	if txType != "" {
		qType = `AND "type" = @txType`
	}

	log.Println(qType)

	query := `SELECT "category", SUM("amount") as "total" FROM public."tx" WHERE "user_id" = @user_id ` + qType + ` AND EXTRACT(MONTH FROM "date") = @month AND EXTRACT(YEAR FROM "date") =  @year  GROUP BY "category"`
	log.Println(query)
	args := pgx.NamedArgs{
		"user_id": userId,
		"month":   month,
		"txType":  txType,
		"year":    t.Year(),
	}
	rows, err := db.Query(ctx, query, args)
	if err != nil {
		return []model.TxChart{}, err
	}
	for rows.Next() {
		p := model.TxChart{}
		err := rows.Scan(&p.Category, &p.Total)
		if err != nil {
			return nil, fmt.Errorf("unable to scan row: %w", err)
		}
		txs = append(txs, p)
	}

	// txs, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.TxChart])
	// if err != nil {
	// 	return nil, err
	// }
	return txs, nil
}

func (r *TxRepo) DeleteByID(ctx context.Context, id int64) error {
	db := database.DB

	query := `DELETE FROM public."tx" WHERE "id" = @id`
	args := pgx.NamedArgs{
		"id": id,
	}
	_, err := db.Exec(ctx, query, args)
	if err != nil {
		return err
	}
	return nil
}

func (r *TxRepo) UpdateByID(ctx context.Context, tx model.Tx) error {
	db := database.DB

	now := time.Now()
	query := `UPDATE public."tx" SET "category" = @category, "date" = @date, "amount" = @amount, "type" = @type, "desc" = @desc, "updated_at" = @updated_at WHERE "user_id" = @user_id AND "id" = @id`
	args := pgx.NamedArgs{
		"id":         tx.ID,
		"category":   tx.Category,
		"date":       tx.Date,
		"user_id":    tx.UserID,
		"amount":     tx.Amount,
		"type":       tx.Type,
		"desc":       tx.Desc,
		"updated_at": now,
	}
	_, err := db.Exec(ctx, query, args)
	if err != nil {
		return err
	}
	return nil
}
