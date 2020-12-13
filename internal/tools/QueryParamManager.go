package tools

import (
	"fmt"
	"strconv"

	"github.com/labstack/echo/v4"
)

type QPM struct {
	Limit  int
	Since  string
	Desc   bool
	Sort   string
	layout string
}

func NewQPM(ctx echo.Context) *QPM {
	limit, _ := strconv.Atoi(ctx.QueryParam("limit"))
	since := ctx.QueryParam("since")
	desc := ctx.QueryParam("desc")
	sort := ctx.QueryParam("sort")

	qpm := &QPM{}

	if desc == "true" {
		qpm.Desc = true
	} else {
		qpm.Desc = false
	}
	qpm.Limit = limit
	qpm.Since = since
	qpm.Sort = sort
	qpm.layout = "2006-01-02T15:04:05.000Z"
	return qpm
}

func (qpm *QPM) UpdateThreadQuery(query string) string {
	if qpm.Desc {
		if qpm.Since != "" {
			query += fmt.Sprintf(" and threads.created <= '%s' ", qpm.Since)
		}
		query += fmt.Sprintf(" order by threads.created desc ")
	} else {
		if qpm.Since != "" {
			query += fmt.Sprintf(" and threads.created >= '%s' ", qpm.Since)
		}
		query += fmt.Sprintf(" order by threads.created")
	}

	if qpm.Limit > 0 {
		query += fmt.Sprintf(" limit %d ", qpm.Limit)
	}

	return query
}
