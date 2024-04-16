package urlProcess

import (
	"main/pkg/storage"
)

func GetNewLink(url string) (string, error) {
	newLink, err := GenerateUniqueCode()
	if err != nil {
		return "", err
	}

	existCode, err := storage.LinkManager(url, newLink, "save")

	if existCode != "" {
		return existCode, nil
	}
	return newLink, nil
}
