package song

import (
	"context"
	"fmt"
	// "strings"
	"errors"
	"musiclib/models/song"
	postgresql "musiclib/pkg/client"

	sq "github.com/Masterminds/squirrel"
)

type repository struct {
	client postgresql.Client
	psql   sq.StatementBuilderType
}

func applyDateFilter(query sq.SelectBuilder, filter *song.ReleaseDateFilter) error {
	switch filter.Condition {
	case "lt":
		query = query.Where(sq.Lt{"release_date": filter.From})
	case "gt":
		query = query.Where(sq.Gt{"release_date": filter.From})
	case "eq":
		query = query.Where(sq.Eq{"release_date": filter.From})
	case "between":
		query = query.Where(sq.GtOrEq{"release_date": filter.From}).
			Where(sq.LtOrEq{"release_date": filter.To})
	default:
		return errors.New("unknown date filter condition")
	}

	return nil
}

func (r *repository) Create(ctx context.Context, s *song.Song) error {
	q := `INSERT INTO tender 
			(name, group, text, link, release_date) 
		  VALUES
			($1, $2, $3, $4, $5)
		  RETURNING id, status, created_at, version`
	err := r.client.QueryRow(
		ctx, q, s.Name, s.Group,
		s.Text, s.Link, s.ReleaseDate).Scan(&s.ID, &s.CreatedAt)
	return err
}

func (r *repository) FindAll(ctx context.Context, qm *song.FindAllQueryModifier) ([]song.Song, error) {
	q := r.psql.Select("id", "name", "group_name", "text", "link", "release_date").From("song").Where(sq.Eq{"is_deleted": "false"})

	if qm.Name != "" {
		q = q.Where(sq.Eq{"name": qm.Name})
	}
	if qm.Group != "" {
		q = q.Where(sq.Eq{"group_name": qm.Group})
	}
	if qm.Text != "" {
		q = q.Where(sq.Eq{"text": qm.Text})
	}
	if qm.Link != "" {
		q = q.Where(sq.Eq{"link": qm.Link})
	}

	if qm.Date != nil {
		err := applyDateFilter(q, qm.Date)
		if err != nil {
			return nil, err
		}
	}
	q.Offset(uint64(qm.Pagination.Offset)).Limit(uint64(qm.Pagination.Limit))
	sql, args, err := q.ToSql()
	fmt.Println(sql)
	if err != nil {
		return nil, err
	}

	rows, err := r.client.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var songs []song.Song
	for rows.Next() {
		var song song.Song
		err := rows.Scan(&song.ID, &song.Name, &song.Group, &song.Text, &song.Link, &song.ReleaseDate)
		if err != nil {
			return nil, err
		}
		songs = append(songs, song)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return songs, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	qb := r.psql.Update("song").Where(sq.Eq{"id": id}).Set("is_deleted", "true")
	sql, args, err := qb.ToSql()
	if err != nil {
		return err
	}

	_, err = r.client.Exec(ctx, sql, args...)
	return err
}

func (r *repository) FindOne(ctx context.Context, id string) (song.Song, error) {
	q := `SELECT 
			id, name, group_name, text, link, release_date, created_at
		  FROM song
		  WHERE id = $1
		  `
	row := r.client.QueryRow(ctx, q, id)
	var s song.Song
	err := row.Scan(&s.ID, &s.Name, &s.Group, &s.Text, &s.Link, &s.ReleaseDate, &s.CreatedAt)
	if err != nil {
		return song.Song{}, err
	}
	return s, nil
}

func (r *repository) Update(ctx context.Context, s *song.Song, fields map[string]interface{}) error {
	qb := r.psql.Update("song").Where(sq.Eq{"id": s.ID}).SetMap(fields).Suffix("RETURNING id, name, group_name, text, link, release_date, created_at")
	sql, args, err := qb.ToSql()
	if err != nil {
		return err
	}
	err = r.client.QueryRow(ctx, sql, args...).Scan(&s.ID, &s.Name, &s.Group, &s.Text, &s.Link, &s.ReleaseDate, &s.CreatedAt)
	return err
}

func NewRepository(clinet postgresql.Client) song.Repository {
	return &repository{
		client: clinet,
		psql:   sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}
