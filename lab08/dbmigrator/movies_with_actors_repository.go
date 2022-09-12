package dbmigrator

import (
	"context"
	"database/sql"
	"lab08/domain"
	"fmt"
	"strings"

	"github.com/alecthomas/log4go"
)

const (
	createTableStmt = `CREATE TABLE IF NOT EXISTS movies_with_actors(
							 id SERIAL PRIMARY KEY,
							 movie_id char(36),
							 actor_id char(36),
							 CONSTRAINT fk_movie
								FOREIGN KEY(movie_id)
									REFERENCES movies(id),
							 CONSTRAINT fk_actor
								 FOREIGN KEY(actor_id)
									 REFERENCES actors(id)
							);
						`
	selectInnerJoinStmt = `SELECT m.id, m.name, m.description, a.id, a.name, a.quote 
								from MOVIES AS m INNER JOIN movies_with_actors AS mwa on m.id=mwa.movie_id 
									INNER JOIN actors AS a ON mwa.actor_id=a.id 
							ORDER BY mwa.id
							LIMIT $1 OFFSET $2`
	insertBulkStmt = "INSERT INTO movies_with_actors (movie_id, actor_id) VALUES %s"
)

type movieWithActorRow struct {
	// TODO: task 3.b.: fill all the necesary fields to parse a row from the JOIN statement
	// movieID
	// movieName
	// movie
}

type MoviesWithActorsRepository struct {
	db        *sql.DB
	movies    domain.Movies

	logger    *log4go.Logger
	batchSize int
	
}

func NewMoviesWithActorsRepository(
	db *sql.DB,
	movies domain.Movies,
	logger *log4go.Logger,
	batchSize int) *MoviesWithActorsRepository {

	mwa := MoviesWithActorsRepository{
		db:        db,
		movies:    movies,
		logger:    logger,
		batchSize: batchSize,
	}

	return &mwa
}

func (mwa *MoviesWithActorsRepository) MigrateMovies(ctx context.Context) error {
	// TODO: task 3.a. create the table movies_with_actors

	// TODO: task 3.a. do a full table scan over the movies table
	// hint: use the mwa.movies.ListMovies() call and fixed limit size
	// the stop condition should be, when we don't have any more records
	
	// TODO: task 3.a. retrieve a batch of records from the movies table

	// TODO: task 3.a. create batch insert statements for each pair ('movie_id', 'actor_id')

	// TODO: task 3.a. append all stored pairs ('movie_id', 'actor_id') (comma separated) and perform the operation
	
	_, err := mwa.db.ExecContext(ctx, createTableStmt)
	if err != nil {
		return err
	}

	for offset := 0; ; offset += mwa.batchSize {
		movies, err := mwa.movies.ListMovies(ctx, mwa.batchSize, offset)
		if err != nil {
			return err
		}
		if len(movies) == 0 {
			return nil
		}

		var valueArgs []string
		for _, movie := range movies {
			for _, actorID := range movie.Actors {
				valueArgs = append(valueArgs, fmt.Sprintf("('%s', '%s')", movie.ID, actorID))
			}
		}

		stmt := fmt.Sprintf(insertBulkStmt, strings.Join(valueArgs, ","))
		_, err = mwa.db.ExecContext(ctx, stmt)
		if err != nil {
			return err
		}
	}
}

func (mwa *MoviesWithActorsRepository) ListAllMovies(ctx context.Context, limit, offset int) ([]domain.MovieExt, error) {
	// TODO: task 3.b.: run the select statement with the inner joins.

	// TODO: task 3.b.: use an intermediary struct to read all the rows

	// TODO: task 3.b.: emit the final array of objects
	rows, err := mwa.db.QueryContext(ctx, selectInnerJoinStmt, limit, offset)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	defer func() {
		closeErr := rows.Close()
		if closeErr != nil {
			_ = mwa.logger.Error(err, domain.FailedRowsErrMsg)
		}
	}()

	actorsWithMovies := make(map[string]map[string]movieWithActorRow)
	for rows.Next() {
		var tmpRow movieWithActorRow
		if err := rows.Scan(&tmpRow.movieID, &tmpRow.movieName, &tmpRow.movieDescription, &tmpRow.actorID, &tmpRow.actorName, &tmpRow.actorQuote); err != nil {
			return nil, err
		}

		actors, found := actorsWithMovies[tmpRow.movieID]
		if found {
			_, found = actors[tmpRow.actorID]
			if !found {
				actors[tmpRow.actorID] = tmpRow
			}
		} else {
			actorsWithMovies[tmpRow.movieID] = make(map[string]movieWithActorRow)
			actorsWithMovies[tmpRow.movieID][tmpRow.actorID] = tmpRow
		}
	}

	var movies []domain.MovieExt
	for _, actorsMap := range actorsWithMovies {
		var m domain.MovieExt
		for _, row := range actorsMap {
			a := domain.Actor{
				Name:  row.actorName,
				Quote: row.actorQuote,
			}
			m.Name = row.movieName
			m.Description = row.movieDescription
			m.Actors = append(m.Actors, a)
		}

		movies = append(movies, m)
	}

	return movies, nil
}
