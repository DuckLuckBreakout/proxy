package proxy

import (
	"github.com/DuckLuckBreakout/proxy/internal/pkg/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UseCase interface {
	HandleHttpRequest(writer http.ResponseWriter, interceptedHttpRequest *http.Request) (string, error)
	HandleHttpsRequest(writer http.ResponseWriter, interceptedHttpRequest *http.Request, needSave bool) error
	DoHttpRequest(httpRequest *http.Request) (*http.Response, error)

	SaveReqToDB(request *http.Request, scheme string, params *gin.Params) error
	GetRequest(id int64) (*models.Request, error)
	GetAllRequests() ([]*models.Request, error)
}
