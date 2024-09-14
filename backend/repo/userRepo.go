package repo

import (
	"context"
	"time"

	"github.com/haqqer/keuanganku/database"
	"github.com/haqqer/keuanganku/model"
	"github.com/jackc/pgx/v5"
)

type UserRepo struct{}

func (r *UserRepo) Create(ctx context.Context, user model.User) error {
	db := database.DB

	now := time.Now()
	query := `INSERT INTO public."users" ("email","username", "name", "google_id", "picture_url", "created_at", "updated_at") VALUES (@email, @username, @name, @google_id, @picture_url, @created_at, @updated_at)`
	args := pgx.NamedArgs{
		"email":       user.Email,
		"username":    user.Username,
		"name":        user.Name,
		"google_id":   user.GoogleID,
		"picture_url": user.PictureUrl,
		"created_at":  now,
		"updated_at":  now,
	}
	_, err := db.Exec(ctx, query, args)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepo) GetByGoogleID(ctx context.Context, googleId string) (model.User, error) {
	db := database.DB
	var user model.User

	query := `SELECT * FROM public."users" WHERE "google_id" = @google_id`
	args := pgx.NamedArgs{
		"google_id": googleId,
	}
	err := db.QueryRow(ctx, query, args).Scan(&user.ID, &user.Username, &user.Name, &user.Email, &user.GoogleID, &user.PictureUrl, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return user, nil
		}
		return user, err
	}
	return user, nil
}
