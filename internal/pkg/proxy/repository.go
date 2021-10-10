package proxy

import "github.com/DuckLuckBreakout/proxy/internal/pkg/models"

type Repository interface {
	InsertRequest(request *models.Request) (int64, error)
	InsertResponse(response *models.Response) error
	GetRequest(id int64) (*models.Request, error)
	GetAllRequests() ([]*models.Request, error)
}
