package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math/rand"
	"time"

	"github.com/sirupsen/logrus"
	"lab09/domain"
	"lab09/gateways/repositories"
)

const (
	localRedisHost   = ""
	defaultRedisPort = 6379
)
func makeSha(b domain.BookMetadata) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(b.Title + b.Author)))
}

func main() {
	log := logrus.New()
	ctx := context.Background()

	redisURL := fmt.Sprintf("%s:%d", localRedisHost, defaultRedisPort)
	redisRepo := repositories.NewRedisRepoFromURL(redisURL, 0)

	rand.Seed(time.Now().Unix())
	timestamp := time.Now().Format(time.RFC3339Nano)
	randomSHA := fmt.Sprintf("%x", sha256.Sum256([]byte(timestamp)))

	metadata := &domain.BookMetadata{
		Sha256:          randomSHA,
		Title:        	 "After",
		Author: 		 "FF",
	}

	metadata2 := &domain.BookMetadata{
		Sha256:          randomSHA,
		Title:        	 "Twilight",
		Author: 		 "FF",
	}

	metadata.Sha256 = makeSha(*metadata)
	metadata2.Sha256 = makeSha(*metadata2)

	
	err := redisRepo.SaveBook(ctx, metadata)
	if err != nil {
		log.WithError(err).Fatal("Could not save data")
	}
	log.Infof("Successfully added metadata %+v", metadata)

	metadata, err = redisRepo.RetrieveBook(ctx, metadata.Sha256)
	if err != nil {
		log.WithError(err).Fatal("Could not load data")
	}
	log.Infof("Successfully loaded metadata %+v", metadata)
	metadata, err = redisRepo.RetrieveBook(ctx, "this_does_not_exist")
	if err != nil {
		log.WithError(err).Fatal("Could not load data")
	}
	if metadata == nil {
		log.Infof("not found")
	}

	err = redisRepo.SaveBook(ctx, metadata2)
	if err != nil {
		log.WithError(err).Fatal("Could not save data")
	}
	log.Infof("Successfully added metadata %+v", metadata2)

	metadata2, err = redisRepo.RetrieveBook(ctx, metadata2.Sha256)
	if err != nil {
		log.WithError(err).Fatal("Could not load data")
	}
	log.Infof("Successfully loaded metadata %+v", metadata2)
	metadata2, err = redisRepo.RetrieveBook(ctx, "this_does_not_exist")
	if err != nil {
		log.WithError(err).Fatal("Could not load data")
	}
	if metadata2 == nil {
		log.Infof("not found")
	}
}
