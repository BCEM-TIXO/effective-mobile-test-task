package song

import (
	"context"
	repeatable "musiclib/pkg/utils"
)

type FindAllQueryModifier struct {
	Pagination repeatable.Pagination
	Name       string
	Group      string
	Text       string
	Link       string
	Date       *ReleaseDateFilter
}

type ReleaseDateFilter struct {
	Condition string
	From      string
	To        string
}

type Repository interface {
	Create(ctx context.Context, t *Song) error
	FindAll(ctx context.Context, qm *FindAllQueryModifier) ([]Song, error)
	FindOne(ctx context.Context, id string) (Song, error)
	Update(ctx context.Context, t *Song, fields map[string]interface{}) error
	Delete(ctx context.Context, id string) error
}
