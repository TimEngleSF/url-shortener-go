package main

import "github.com/TimEngleSF/url-shortener-go/internal/models"

type templateData struct {
	CurrentYear int
	Link        *models.Link
	Validation  map[string]string
}
