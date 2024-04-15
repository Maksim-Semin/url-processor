package storage

import (
	"fmt"
	"go.uber.org/zap"
	lg "main/pkg/logger"
	"main/pkg/storage/imdb"
	"main/pkg/storage/postgres"
	"main/resource"
)

func LinkManager(sourceURL, shortLink, action string) (string, error) {
	if resource.CFG.DB == "postgres" {
		if action == "save" {
			shortCode, err := postgres.Db.SaveUrl(shortLink, sourceURL)
			if err != nil {
				lg.Logger.Warn("error on saving in PostgreSQL", zap.Error(err))
			}
			if shortCode != "" {
				return shortCode, nil
			}
		} else {
			sourceURL, err := postgres.Db.GetURL(shortLink)
			if err != nil {
				lg.Logger.Warn("error on taking link from PostgreSQL", zap.Error(err))
				return "", fmt.Errorf("error on taking link from DB")
			}
			return sourceURL, nil
		}
	} else {
		if action == "save" {
			imdb.MemoryDB.Set(shortLink, sourceURL)
		} else {
			sourceURL, ok := imdb.MemoryDB.Get(shortLink)
			if !ok {
				lg.Logger.Warn("error on taking link from map")
				return "", fmt.Errorf("error on taking link from DB")
			}
			return sourceURL, nil
		}

	}
	return "", nil
}
