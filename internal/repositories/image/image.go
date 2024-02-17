package image

import (
	"context"
	"log/slog"

	"github.com/ijimiji/pipeline/internal/models"
	"github.com/ijimiji/pipeline/internal/services/sqlite"
)

func New(db *sqlite.Database) *Repository {
	if _, err := db.ExecContext(context.Background(), `create table if not exists Images (image_id TEXT, prompt TEXT, url TEXT, status TEXT);`); err != nil {
		panic(err)
	}
	return &Repository{
		db: db,
	}
}

type Repository struct {
	db *sqlite.Database
}

func (r *Repository) Add(ctx context.Context, image models.Image) error {
	_, err := r.db.ExecContext(
		ctx,
		`insert into Images (image_id, prompt, url, status) values ($1, $2, $3, $4)`,
		image.ID,
		image.Prompt,
		image.URL,
		image.GenerationStatus,
	)
	if err != nil {
		slog.Error(err.Error())
	}

	return err
}

func (r *Repository) Get(ctx context.Context, id string) (models.Image, error) {
	var ret models.Image
	row := r.db.QueryRowContext(
		ctx,
		`select image_id, prompt, url, status from Images where image_id = $1 limit 1`,
		id,
	)

	if err := row.Scan(
		&ret.ID,
		&ret.Prompt,
		&ret.URL,
		&ret.GenerationStatus,
	); err != nil {
		slog.Error(err.Error())
		return ret, err
	}

	return ret, nil
}

func (r *Repository) SetURL(ctx context.Context, id string, url string) error {
	_, err := r.db.ExecContext(
		ctx,
		`
		update Images
		set 
			url = $1, 
			status = "ready"
		where image_id = $2
		`,
		url,
		id,
	)
	if err != nil {
		slog.Error(err.Error())
	}

	return err
}
