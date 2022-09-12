package sql

import (
	"context"
	"database/sql"
	// "fmt"
	"lab08/domain"
	// "strings"

	"github.com/alecthomas/log4go"
	"github.com/lib/pq"
)

const (
	sqlActorsCreateStmt = `INSERT INTO actors (id, name, quote)
					VALUES ($1, $2, $3) 
					RETURNING id, name, quote`
	sqlActorsDeleteStmt = `DELETE FROM actors WHERE id = $1 
					RETURNING id, name, quote`
	sqlActorsGetByIDStmt = `SELECT id, name, quote
					FROM actors 
					WHERE id = ANY ($1)`
	listActorsStmt = `SELECT id, name, quote
				FROM actors 
				ORDER BY id
				LIMIT $1 OFFSET $2`
	insertActorsBulkStmt        = "INSERT INTO actors (id, name, quote) VALUES %s"
	countActorsPartialMatchStmt = "SELECT count(*) from actors where name like"
	countActorsExactMatchStmt   = "SELECT count(*) from actors where name="
)

type ActorRepository struct {
	db     *sql.DB
	logger *log4go.Logger
}

func NewActorRepository(db *sql.DB, logger *log4go.Logger) *ActorRepository {
	mr := ActorRepository{
		db:     db,
		logger: logger,
	}

	return &mr
}

func (mr *ActorRepository) CreateActor(ctx context.Context, actor domain.Actor) error {
	// TODO: implement me
	err := mr.db.QueryRowContext(ctx, sqlActorsCreateStmt, actor.ID, actor.Name, actor.Quote)
	if err != nil && err.Err() != nil {
		return err.Err()
	}
	return nil
}

func (mr *ActorRepository) CreateActors(ctx context.Context, actors []domain.Actor) error {
	// TODO: task 5: index bulk data in the actors table
	return nil
}

func (mr *ActorRepository) GetActorsByID(ctx context.Context, ids []string) ([]domain.Actor, error) {
	// TODO: implement me
	rows, err := mr.db.QueryContext(ctx, sqlActorsGetByIDStmt, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	defer func() {
		closeErr := rows.Close()
		if closeErr != nil {
			_ = mr.logger.Error(err, domain.FailedRowsErrMsg)
		}
	}()

	var actors []domain.Actor
	for rows.Next() {
		actor := domain.Actor{}

		if err := rows.Scan(&actor.ID, &actor.Name, &actor.Quote); err != nil {
			return nil, err
		}

		actors = append(actors, actor)
	}

	return actors, nil
}

func (mr *ActorRepository) ListActors(ctx context.Context, limit, offset int) ([]domain.Actor, error) {
	// TODO: implement me

	rows, err := mr.db.QueryContext(ctx, listActorsStmt, limit, offset)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	defer func() {
		closeErr := rows.Close()
		if closeErr != nil {
			_ = mr.logger.Error(err, domain.FailedRowsErrMsg)
		}
	}()

	var actors []domain.Actor
	for rows.Next() {
		var actor domain.Actor
		if err := rows.Scan(&actor.ID, &actor.Name, &actor.Quote); err != nil {
			return nil, err
		}
		actors = append(actors, actor)
	}

	return actors, nil
}

func (mr *ActorRepository) DeleteActor(ctx context.Context, id string) error {
	// TODO: implement me
	rows, err := mr.db.QueryContext(ctx, sqlActorsDeleteStmt, id)
	if err != nil {
		return err
	}
	if rows.Err() != nil {
		return rows.Err()
	}
	defer func() {
		closeErr := rows.Close()
		if closeErr != nil {
			_ = mr.logger.Error(err, domain.FailedRowsErrMsg)
		}
	}()

	return nil
}

func (mr *ActorRepository) Count(ctx context.Context, name string, useExactMatch bool) (int, error) {
	var stmt string
	if useExactMatch {
		stmt = countActorsExactMatchStmt + "'" + name + "'"
	} else {
		stmt = countActorsPartialMatchStmt + "'%" + name + "%'"
	}

	// TODO: also remove this log line, as it will get quite spammy
	mr.logger.Info("Running count statement: %s", stmt)
	// TODO: task 5.b.:

	return 0, nil
}