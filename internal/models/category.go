package models

import "git.ramadhantriyant.id/ramadhantriyant/substrack/internal/database"

type CategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CategoryList struct {
	Total      int                 `json:"total"`
	Categories []database.Category `json:"categories"`
}
