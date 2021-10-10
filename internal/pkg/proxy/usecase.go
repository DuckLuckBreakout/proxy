package proxy

import (
	"github.com/DuckLuckBreakout/proxy/internal/pkg/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UseCase interface {
	HandleHttpRequest(writer http.ResponseWriter, interceptedHttpRequest *http.Request, requestId int64) (string, error)
	HandleHttpsRequest(writer http.ResponseWriter, interceptedHttpRequest *http.Request, needSave bool) error
	DoHttpRequest(httpRequest *http.Request) (*http.Response, error)

	SaveReqToDB(request *http.Request, scheme string, params *gin.Params) (int64, error)
	GetRequest(id int64) (*models.Request, error)
	GetAllRequests() ([]*models.Request, error)
}
