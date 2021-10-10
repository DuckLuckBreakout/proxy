package models

type Request struct {
	Id      int64
	Method  string
	Scheme  string
	Host    string
	Path    string
	Headers map[string][]string
	Body    string
	Params map[string]string
}

type Response struct {
	Id        int64
	RequestId int64
	Status      int
	Headers   map[string][]string
	Body      string
}