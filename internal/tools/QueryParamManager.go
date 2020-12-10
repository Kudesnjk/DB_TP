package tools

import (
	"fmt"
	"strconv"

	"github.com/labstack/echo/v4"
)

type QPM struct {
	Limit int
	Since string
	Desc  bool
	Sort  string
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
	return qpm
}

func (qpm *QPM) UpdateThreadQuery(query string) string {
	if qpm.Since != "" {
		query += fmt.Sprintf(" and threads.created >= %s ", qpm.Since)
	}

	if qpm.Desc {
		query += fmt.Sprintf(" order by threads.created ")
	} else {
		query += fmt.Sprintf(" order by threads.created desc")
	}

	if qpm.Limit > 0 {
		query += fmt.Sprintf(" limit %d ", qpm.Limit)
	}
	return query
}
