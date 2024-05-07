package cache

import (
	"encoding/json"
	"fmt"

	"github.com/woaitsAryan/stuneckt-task/models"
)

func GetPost(id string) (*models.Post, error) {
	data, err := GetFromCache("post-" + id)
	if err != nil {
		return nil, err
	}

	post := models.Post{}
	if err = json.Unmarshal([]byte(data), &post); err != nil {
		return nil, fmt.Errorf("error while unmarshaling post: %w", err)
	}
	return &post, nil
}

func SetPost(id string, post *models.Post) error {
	data, err := json.Marshal(post)
	if err != nil {
		return fmt.Errorf("error while marshaling post: %w", err)
	}
	if err := SetToCache("post-"+id, data); err != nil {
		return err
	}
	return nil
}

func RemovePost(id string) error {
	if err := RemoveFromCache("post-" + id); err != nil {
		return err
	}
	return nil
}