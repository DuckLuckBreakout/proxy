package usecase

import (
	"compress/gzip"
	"encoding/json"
	"github.com/DuckLuckBreakout/proxy/internal/pkg/proxy"
	https "github.com/DuckLuckBreakout/proxy/internal/pkg/https"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/DuckLuckBreakout/proxy/internal/pkg/models"
)

type ProxyUseCase struct {
	proxyRepo proxy.Repository
}

func NewUseCase(proxyRepo proxy.Repository) proxy.UseCase {
	return &ProxyUseCase{
		proxyRepo: proxyRepo,
	}
}

func (u *ProxyUseCase) SaveReqToDB(request *http.Request, scheme string, params *gin.Params) error {
	bodyBytes, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return err
	}

	var allParams map[string]string
	json.Unmarshal(bodyBytes, &allParams)
	if allParams == nil {
		allParams = make(map[string]string)
	}
	if params != nil {
		for _, val := range *params {
			allParams[val.Key] = val.Value
		}
	}


	req := &models.Request{
		Method:  request.Method,
		Scheme:  scheme,
		Host:    request.Host,
		Path:    request.URL.Path,
		Headers: request.Header,
		Body:    string(bodyBytes),
		Params: allParams,
	}

	err = u.proxyRepo.InsertInto(req)
	if err != nil {
		return err
	}

	return nil
}

func (u *ProxyUseCase) GetRequest(id int64) (*models.Request, error) {
	return u.proxyRepo.GetRequest(id)
}

func (u *ProxyUseCase) GetAllRequests() ([]*models.Request, error) {
	return u.proxyRepo.GetAllRequests()
}

func decodeResponse(response *http.Response) ([]byte, error) {
	var body io.ReadCloser

	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		var err error
		body, err = gzip.NewReader(response.Body)
		if err != nil {
			body = response.Body
		}
	default:
		body = response.Body
	}

	bodyByte, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	lineBreak := []byte("\n")
	bodyByte = append(bodyByte, lineBreak...)
	bodyByte = bodyByte[0:500]

	defer body.Close()

	return bodyByte, nil
}

func (u *ProxyUseCase) HandleHttpRequest(writer http.ResponseWriter, interceptedHttpRequest *http.Request) (string, error) {
	proxyResponse, err := u.DoHttpRequest(interceptedHttpRequest)
	if err != nil {
		panic(err)
	}
	for header, values := range proxyResponse.Header {
		for _, value := range values {
			writer.Header().Add(header, value)
		}
	}
	writer.WriteHeader(proxyResponse.StatusCode)
	_, err = io.Copy(writer, proxyResponse.Body)
	if err != nil {
		panic(err)
	}

	decodedResponse, err := decodeResponse(proxyResponse)
	if err != nil {
		return "", err
	}

	defer proxyResponse.Body.Close()
	return string(decodedResponse), nil
}

func (u *ProxyUseCase) HandleHttpsRequest(writer http.ResponseWriter, interceptedHttpRequest *http.Request, needSave bool) error {
	httpsService := https.NewHttpsService(writer, interceptedHttpRequest, u.proxyRepo)

	err := httpsService.ProxyHttpsRequest()
	if err != nil {
		return err
	}
	return nil
}

func (u *ProxyUseCase) DoHttpRequest(HttpRequest *http.Request) (*http.Response, error) {
	proxyRequest, err := http.NewRequest(HttpRequest.Method, HttpRequest.URL.String(), HttpRequest.Body)
	if err != nil {
		return nil, err
	}

	proxyRequest.Header = HttpRequest.Header

	proxyResponse, err := http.DefaultTransport.RoundTrip(proxyRequest)
	if err != nil {
		return nil, err
	}

	return proxyResponse, nil
}
