package main

import "encoding/json"
import "fmt"

type query struct {
	req string
	res interface{}
	err error
}

func newQuery(req string) *query {
	q := query{
		req: req,
		res: nil,
		err: nil,
	}
	return &q
}

func (q *query) parse(timeout uint16) {
	q.res, q.err = client.Query(q.req, timeout)
}

func (q *query) json() (string, error) {
	var b []byte
	var err error
	if q.res == nil {
		return "", fmt.Errorf("nothing to JSONify")
	}
	if b, err = json.Marshal(q.res); err != nil {
		return "", err
	}
	return string(b), nil
}
