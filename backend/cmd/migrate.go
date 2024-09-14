package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/haqqer/keuanganku/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Migrate() {
	ctx := context.Background()
	var argsRaw = os.Args
	if len(argsRaw) <= 2 {
		fmt.Println("use migrate `up` or `down`")
		os.Exit(1)
	}
	arg := argsRaw[2]
	pool, err := pgxpool.NewWithConfig(ctx, database.Config())

	if err != nil {
		log.Fatalf("error db connection")
	}

	switch arg {
	case "up":
		up(ctx, pool)
	case "down":
		down(ctx, pool)
	}
}

func up(ctx context.Context, db *pgxpool.Pool) {
	queryFile, err := os.ReadFile("./database/up.sql")
	if err != nil {
		log.Fatal("error read file")
	}
	query := string(queryFile)
	_, err = db.Exec(ctx, query)
	if err != nil {
		log.Fatalf("error execute query : %s", err.Error())
	}
	log.Println("migration up success")
}

func down(ctx context.Context, db *pgxpool.Pool) {
	queryFile, err := os.ReadFile("./database/down.sql")
	if err != nil {
		log.Fatal("error read file")
	}
	query := string(queryFile)
	_, err = db.Exec(ctx, query)
	if err != nil {
		log.Fatalf("error execute query : %s", err.Error())
	}
	log.Println("migration down success")
}
