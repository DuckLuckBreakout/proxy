package proxy

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler interface {
	ScanRequest(c *gin.Context, params map[string]string)
	ScanRequestHandler(c *gin.Context)
	RepeatRequestHandler(c *gin.Context)
	GetAllRequestsHandler(c *gin.Context)
	GetRequestHandler(c *gin.Context)
	HandleRequest(c *gin.Context)
	HandleRequestHttp(writer http.ResponseWriter, request *http.Request)
}
