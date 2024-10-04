package song

import (
	"context"
	// "strings"
	// "errors"
	"musiclib/models/song"
	postgresql "musiclib/pkg/client"

	sq "github.com/Masterminds/squirrel"
)

type repository struct {
	client postgresql.Client
	psql   sq.StatementBuilderType
}

func (r *repository) Create(ctx context.Context, t *song.Song) error {
	panic("unimplemented")
}

func (r *repository) FindAll(ctx context.Context, qm *song.FindAllQueryModifier) ([]song.Song, error) {
	panic("unimplemented")
}

func (r *repository) FindOne(ctx context.Context, id string) (song.Song, error) {
	panic("unimplemented")
}

func (r *repository) Update(ctx context.Context, t *song.Song, fields map[string]interface{}) error {
	panic("unimplemented")
}

func NewRepository(clinet postgresql.Client) song.Repository {
	return &repository{
		client: clinet,
		psql:   sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}
