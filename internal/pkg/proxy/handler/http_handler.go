package handler

import (
	"github.com/DuckLuckBreakout/proxy/internal/pkg/proxy"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	needSave   bool
	proxyUCase proxy.UseCase
}

func (h *UserHandler) HandleRequestHttp(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodConnect {
		err := h.proxyUCase.HandleHttpsRequest(writer, request, h.needSave)
		if err != nil {
			panic(err)
		}
	} else {
		var requestId int64
		if h.needSave {
			params := make(gin.Params, 0)
			data, _ := url.ParseQuery(request.URL.RawQuery)
			for key, val := range data {
				params = append(params, gin.Param{
					Key:   key,
					Value: val[0],
				})
			}
			var err error
			requestId, err = h.proxyUCase.SaveReqToDB(request, "http", &params)
			if err != nil {
				panic(err)
			}
		}

		_, err := h.proxyUCase.HandleHttpRequest(writer, request, requestId)
		if err != nil {
			panic(err)
		}
	}
	h.needSave = true
}

func (h *UserHandler) ScanRequest(c *gin.Context, params map[string]string) {
	queryParams := url.Values{}
	for key, val := range params {
		queryParams.Add(key, val + "vulnerable'\"><img src onerror=alert()>")
	}

	c.Request.URL.RawQuery = queryParams.Encode()

	if c.Request.Method == http.MethodPost {
		bodyBytes, err := ioutil.ReadAll(c.Request.Body)

		c.Request.Body = ioutil.NopCloser(strings.NewReader(string(bodyBytes) + "vulnerable'\"><img src onerror=alert()>"))
		if err != nil {
			panic(err)
		}
	}

	response, err := h.proxyUCase.HandleHttpRequest(c.Writer, c.Request, 0)
	if err != nil {
		panic(err)
	}

	if strings.Contains(response, "vulnerable'\"><img src onerror=alert()>") {
		c.Status(http.StatusConflict)
		if err != nil {
			panic(err)
		}
		return
	}

	//result := make(map[string]string)
	//result["result"] = "XSS - \n"
	//bytes, _ := json.Marshal("XSS")
	c.Status(http.StatusOK)
	return
	//_, err = io.Copy(c.Writer, strings.NewReader("Request not contains XSS injection\n"))
	//if err != nil {
	//	panic(err)
	//}
}

func (h *UserHandler) ScanRequestHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	foundReq, err := h.proxyUCase.GetRequest(id)
	if err != nil {
		panic(err)
	}

	c.Request = &http.Request{
		Method: foundReq.Method,
		URL: &url.URL{
			Scheme: foundReq.Scheme,
			Host:   foundReq.Host,
			Path:   foundReq.Path,
		},
		Header: foundReq.Headers,
		Body:   ioutil.NopCloser(strings.NewReader(foundReq.Body)),
		Host:   c.Request.Host,
	}
	h.ScanRequest(c, foundReq.Params)
	return
	//if err != nil {
	//	panic(err)
	//}
}

func (h *UserHandler) RepeatRequestHandler(c *gin.Context) {
	h.needSave = false
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	foundReq, err := h.proxyUCase.GetRequest(id)
	if err != nil {
		panic(err)
	}

	c.Request = &http.Request{
		Method: foundReq.Method,
		URL: &url.URL{
			Scheme: foundReq.Scheme,
			Host:   foundReq.Host,
			Path:   foundReq.Path,
		},
		Header: foundReq.Headers,
		Body:   ioutil.NopCloser(strings.NewReader(foundReq.Body)),
		Host:   c.Request.Host,
	}

	h.HandleRequest(c)
}

func (h *UserHandler) GetAllRequestsHandler(c *gin.Context) {
	requests, err := h.proxyUCase.GetAllRequests()
	if err != nil {
		panic(err)
	}

	var result string
	for _, request := range requests {
		result += "\nId: " + strconv.FormatInt(request.Id, 16) + "\nMethod: " + request.Method +
			"\nScheme: " + request.Scheme + "\nPath: " + request.Path + "\nHost: " + request.Host +
			"\nBody: " + request.Body + "\n" + "Params: \n"
		for key, param := range request.Params {
			result += "\t" + key + "="+ param + "\n"
		}
	}
	_, err = io.Copy(c.Writer, strings.NewReader(result))
	if err != nil {
		panic(err)
	}
}

func (h *UserHandler) GetRequestHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	foundReq, err := h.proxyUCase.GetRequest(id)
	if err != nil {
		panic(err)
	}

	result := "Id: " + strconv.FormatInt(foundReq.Id, 16) + "\nMethod: " + foundReq.Method +
		"\nScheme: " + foundReq.Scheme + "\nPath: " + foundReq.Path + "\nHost: " + foundReq.Host +
		"\nBody: " + foundReq.Body + "\n"

	_, err = io.Copy(c.Writer, strings.NewReader(result))
	if err != nil {
		panic(err)
	}
}

func (h *UserHandler) HandleRequest(c *gin.Context) {
	if c.Request.Method == http.MethodConnect {
		err := h.proxyUCase.HandleHttpsRequest(c.Writer, c.Request, h.needSave)
		if err != nil {
			panic(err)
		}
	} else {
		var requestId int64
		if h.needSave {
			var err error
			requestId, err = h.proxyUCase.SaveReqToDB(c.Request, "http", &c.Params)
			if err != nil {
				panic(err)
			}
		}

		_, err := h.proxyUCase.HandleHttpRequest(c.Writer, c.Request, requestId)
		if err != nil {
			panic(err)
		}
	}
	h.needSave = true
}

func NewHandler(userUCase proxy.UseCase) proxy.Handler {
	return &UserHandler{
		proxyUCase: userUCase,
		needSave:   true,
	}
}

