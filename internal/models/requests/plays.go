package request

import model "theater-ticket-system/internal/models/models"

type Play struct {
	Title       string `json:"title" binding:"required"`
	Author      string `json:"author" binding:"required"`
	Description string `json:"description" binding:"required"`
	Duration    int    `json:"duration" binding:"required"`
	PosterURL   string `json:"poster_url" binding:"required"`
	Genre       string `json:"genre" binding:"required"`
}

func (p *Play) Model() *model.Play {
	return &model.Play{
		Title:       p.Title,
		Author:      p.Author,
		Description: p.Description,
		Duration:    p.Duration,
		PosterURL:   p.PosterURL,
		Genre:       p.Genre,
	}
}
