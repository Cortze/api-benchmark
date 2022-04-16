package requester

import (
	"time"
)

var ()

type Request struct {
	Query        string
	Status       string
	RequestTime  time.Time
	ResponseTime time.Duration
	Error        string
}

func NewRequest(query string, status string, requesttime time.Time, responsetime time.Duration, err string) *Request {
	return &Request{
		Query:        query,
		Status:       status,
		RequestTime:  requesttime,
		ResponseTime: responsetime,
		Error:        err,
	}
}

func RequestStatusCsvColumnNames() string {
	csvline := "query" + "," + "request status" + "," + "request time" + "," + "response time" + "," + "error" + "\n"
	return csvline
}

func (r *Request) CsvLine() string {
	csvline := r.Query + "," + r.Status + "," + r.RequestTime.String() + "," + r.ResponseTime.String() + "\n"
	return csvline
}
