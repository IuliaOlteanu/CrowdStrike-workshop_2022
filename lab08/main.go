package main

import (
	"context"
	"lab08/dbgenerator"
	"lab08/dbmigrator"
	"lab08/domain"
	"lab08/sql"
	"os"
	"strconv"
	"time"

	"github.com/alecthomas/log4go"
	"github.com/montanaflynn/stats"
)

func main() {
	// create context
	ctx := context.Background()

	// create logger
	l := log4go.NewDefaultLogger(log4go.INFO)
	defer l.Close()

	l.Info("Hello world=%+v!", os.Getegid())

	// create DB connection
	db, err := sql.CreatePostgresConnection(
		"0.0.0.0:5432", "",
		"upb",
		"upb",
		"disable")
	if err != nil {
		_ = l.Error("Failed creating connection=%+v", err)
	}

	// use movies moviesRepo
	moviesRepo := sql.NewMovieRepository(db, &l)
	movies, err := moviesRepo.GetMoviesByID(ctx, []string{"001"})
	l.Info("GetMoviesByID: movies=%+v err=%+v", movies, err)

	// TODO: task 2.a. - create 10 movies
	for i := 0; i < 10; i++{
	 	sir := strconv.Itoa(i)
	 	err = moviesRepo.CreateMovie(ctx, domain.Movie{ID:"movie" + sir, Name: "number of movie" + sir, Description: "moviee"})
	 	if err != nil {
	 		l.Error("Failed creating movie:%+v", err)
	 	} 
	}
	err = moviesRepo.CreateMovie(ctx, domain.Movie{ID:"123", Name: "my name isn't given", Description: "that's another story"})
	if err != nil {
		l.Error("Failed creating movie:%+v", err)
	}

	// TODO: task 2.a. - are your movies appearing in the list?
	movies, err = moviesRepo.ListMovies(ctx, 100, 0)
	l.Info("Listmovies: ids=%+v err=%+v", movies, err)

	// TODO: task 2.a. - remove the previously created 10 movies
	for i := 0; i < 10; i++{
		sir := strconv.Itoa(i)
		err = moviesRepo.DeleteMovie(ctx, "movie" + sir)
		if err != nil {
			l.Error("Failed deleting movie:%+v", err)
		} 
	}
	err = moviesRepo.DeleteMovie(ctx, "123")
	if err != nil {
		l.Error("Failed deleting movie:%+v", err)
	} 

	// TODO: task 2.a. - are your movies still showing in the list?
	movies, err = moviesRepo.ListMovies(ctx, 100, 0)
	l.Info("Listmovies: ids=%+v err=%+v", movies, err)

	// TODO: task 2.b. - implement actors repository
	actorsRepo := sql.NewActorRepository(db, &l)
	actors, err := actorsRepo.GetActorsByID(ctx, []string{"001", "002"})
	
	l.Info("GetActorsByID: actors=%+v err=%+v", actors, err)
	actors, err = actorsRepo.ListActors(ctx, 100, 0)
	l.Info("ListActors: actors=%+v err=%+v", actors, err)

	// TODO: task 3: performing join operations
	mwa := dbmigrator.NewMoviesWithActorsRepository(db, moviesRepo, &l, 100)

	// TODO: task 3.a.: perform the movies_with_actors table migration
	err = mwa.MigrateMovies(ctx)
	l.Info("MigrateMovies: err=%+v", err)


	// TODO: task 3.b.: list the full movie names associated with all the actor names
	movieExt, err := mwa.ListAllMovies(ctx, 100, 0)
	l.Info("ListAllMovies: movies=%+v err=%+v", movieExt, err)

	// TODO: task 4: generate bulk random data
	gen := dbgenerator.NewRandomActorGenerator(actorsRepo, &l, 20, 200)
	err = gen.Generate(ctx, 3000000)
	if err != nil {
		l.Info("Generate: err=%+v", err)
	}

	// TODO: task 5: compute the median over several random count(*) from the DB
	var durations []float64
	for i := 0; i < 200; i++ {
		// TODO - can we paralelize this part?
		counter, err := gen.Count(ctx)
		if err != nil {
			l.Info("Count: err=%+v", err)
			continue
		}

		durations = append(durations, float64(counter.ExactMatchDuration))

		// avoid spammy log messages and sample print only some values
		if i%20 == 0 {
			l.Info("Count for for exact match %s=%d (%+v) [partial name=%s count=%d (%+v)]",
				counter.ExactMatchName,
				counter.ExactMatchCount,
				counter.ExactMatchDuration,
				counter.PartialMatchName,
				counter.PartialMatchCount,
				counter.PartialMatchDuration)
		}
	}

	// TODO: task 6: observe what happens when creating an index on the actors table
	durationFloat, _ := stats.Median(durations)
	d := time.Duration(durationFloat)
	l.Info("Median duration is=%v", d)
}
