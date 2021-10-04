package proxy

import "github.com/DuckLuckBreakout/proxy/internal/pkg/models"

type Repository interface {
	InsertInto(request *models.Request) error
	GetRequest(id int64) (*models.Request, error)
	GetAllRequests() ([]*models.Request, error)
}
