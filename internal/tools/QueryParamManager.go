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
	}

	qpm.Limit = limit
	qpm.Since = since
	qpm.Sort = sort
	qpm.layout = "2006-01-02T15:04:05.000Z"

	return qpm
}

func (qpm *QPM) UpdateForumUsersQuery(queryThreads, queryPosts string) string {
	query := ""
	unionStr := " union distinct "

	if qpm.Desc {
		if qpm.Since != "" {
			queryThreads += fmt.Sprintf(" and nickname < '%s' ", qpm.Since)
			queryPosts += fmt.Sprintf(" and nickname < '%s' ", qpm.Since)
		}
		query = queryThreads + unionStr + queryPosts
		query += fmt.Sprintf(" order by nickname desc")
	} else {
		if qpm.Since != "" {
			queryThreads += fmt.Sprintf(" and nickname > '%s' ", qpm.Since)
			queryPosts += fmt.Sprintf(" and nickname > '%s' ", qpm.Since)
		}
		query = queryThreads + unionStr + queryPosts
		query += fmt.Sprintf(" order by nickname ")
	}

	if qpm.Limit > 0 {
		query += fmt.Sprintf(" limit %d ", qpm.Limit)
	}

	return query
}

func (qpm *QPM) UpdatePostQuery(query string) string {
	switch qpm.Sort {
	case "tree":
		if qpm.Desc {
			if qpm.Since != "" {
				query += fmt.Sprintf(" and path < (select path from posts where id = %s) ", qpm.Since)
			}
			query += fmt.Sprintf(" order by posts.path desc, posts.id desc ")
		} else {
			if qpm.Since != "" {
				query += fmt.Sprintf(" and path > (select path from posts where id = %s) ", qpm.Since)
			}
			query += fmt.Sprintf(" order by posts.path, posts.id ")
		}

	case "parent_tree":
		if qpm.Desc {
			if qpm.Since != "" {
				query += fmt.Sprintf(" and path[1] in (select distinct path[1] from posts where path[1] < (select path[1] from posts where id = %s) and array_length(path, 1) = 1 and thread_id = $1 ", qpm.Since)
			} else {
				query += fmt.Sprintf(" and path[1] in (select distinct path[1] from posts where array_length(path, 1) = 1 and thread_id = $1 order by path[1] desc")
			}
			if qpm.Limit > 0 {
				query += fmt.Sprintf(" limit %d ) ", qpm.Limit)
			}
			query += fmt.Sprintf(" order by posts.path[1] desc, posts.path, posts.id ")
		} else {
			if qpm.Since != "" {
				query += fmt.Sprintf(" and path[1] in (select distinct path[1] from posts where path[1] > (select path[1] from posts where id = %s) and array_length(path, 1) = 1 and thread_id = $1 ", qpm.Since)
			} else {
				query += fmt.Sprintf(" and path[1] in (select distinct path[1] from posts where array_length(path, 1) = 1 and thread_id = $1 order by path[1] ")
			}
			if qpm.Limit > 0 {
				query += fmt.Sprintf(" limit %d ) ", qpm.Limit)
			}
			query += fmt.Sprintf(" order by posts.path[1], posts.path, posts.id ")
		}
		return query
	default:
		if qpm.Desc {
			if qpm.Since != "" {
				query += fmt.Sprintf(" and posts.id < %s ", qpm.Since)
			}
			query += fmt.Sprintf(" order by posts.created desc, posts.id desc ")
		} else {
			if qpm.Since != "" {
				query += fmt.Sprintf(" and posts.id > %s ", qpm.Since)
			}
			query += fmt.Sprintf(" order by posts.created, posts.id ")
		}
	}

	if qpm.Limit > 0 {
		query += fmt.Sprintf(" limit %d ", qpm.Limit)
	}

	return query
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
