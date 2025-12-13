package handlers

import "git.ramadhantriyant.id/ramadhantriyant/substrack/internal/models"

type Handler struct {
	config *models.AppConfig
}

func New(config *models.AppConfig) *Handler {
	return &Handler{
		config: config,
	}
}
