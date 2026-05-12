package service

import "github.com/google/uuid"

func generateName() string {
	return uuid.New().String()
}
