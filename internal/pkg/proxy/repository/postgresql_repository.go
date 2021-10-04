package repository

import (
	"encoding/json"
	"github.com/DuckLuckBreakout/proxy/internal/pkg/models"
	"github.com/DuckLuckBreakout/proxy/internal/pkg/proxy"
	"github.com/jmoiron/sqlx"
)

type ProxyRepository struct {
	dbConn *sqlx.DB
}

func (u ProxyRepository) InsertInto(request *models.Request) error {
	bs, err := json.Marshal(request.Headers)
	if err != nil {
		return err
	}

	paramsB, err := json.Marshal(request.Params)
	if err != nil {
		return err
	}

	if err := u.dbConn.QueryRow(
		`INSERT INTO requests(method, scheme, host, path, headers, body, params) 
				VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		request.Method, request.Scheme, request.Host, request.Path, string(bs), request.Body, string(paramsB)).Scan(&request.Id); err != nil {
		return err
	}

	return nil
}

func (u ProxyRepository) GetRequest(id int64) (*models.Request, error) {
	rows, err := u.dbConn.Queryx(`SELECT * from requests where id = $1`, id)
	if err != nil {
		return nil, err
	}

	requests, err := u.scanRequests(rows)
	if err != nil {
		return nil, err
	}

	if len(requests) == 0 {
		return nil, nil
	}

	return requests[0], nil
}

func (u ProxyRepository) GetAllRequests() ([]*models.Request, error) {
	rows, err := u.dbConn.Queryx(`SELECT * from requests`)
	if err != nil {
		return nil, err
	}

	requests, err := u.scanRequests(rows)
	if err != nil {
		return nil, err
	}

	return requests, nil
}

func (u *ProxyRepository) scanRequests(rows *sqlx.Rows) ([]*models.Request, error) {
	var requests []*models.Request
	for rows.Next() {
		reqMap := make(map[string]interface{})
		err := rows.MapScan(reqMap)
		if err != nil {
			return nil, err
		}

		headersRaw := reqMap["headers"].([]byte)

		var headers map[string][]string
		err = json.Unmarshal(headersRaw, &headers)
		if err != nil {
			return nil, err
		}


		paramsRaw := reqMap["params"].([]byte)

		var params map[string]string
		err = json.Unmarshal(paramsRaw, &params)
		if err != nil {
			return nil, err
		}

		requests = append(requests, &models.Request{
			Id:      reqMap["id"].(int64),
			Method:  reqMap["method"].(string),
			Scheme:  reqMap["scheme"].(string),
			Host:    reqMap["host"].(string),
			Path:    reqMap["path"].(string),
			Headers: headers,
			Body:    reqMap["body"].(string),
			Params:  params,
		})
	}

	return requests, nil
}

func NewRepository(conn *sqlx.DB) proxy.Repository {
	return &ProxyRepository{
		dbConn: conn,
	}
}
