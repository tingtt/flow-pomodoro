package pomodoro

import (
	"flow-pomodoro/mysql"
	"time"
)

type GetListQuery struct {
	Start              *time.Time `query:"start" validate:"required"`
	End                *time.Time `query:"end" validate:"required"`
	ProjectId          *uint64    `query:"project_id" validate:"omitempty"`
	IncludeSubProjects bool       `query:"include_sub_project" validate:"omitempty"`
	TodoId             *uint64    `query:"todo_id" validate:"omitempty"`
}

func GetList(userId uint64, q GetListQuery) (pomodoros []Pomodoro, err error) {
	// Generate query
	queryStr := "SELECT id, start, end, remaining_time, todo_id, project_id, parent_project_id FROM logs WHERE user_id = ?"
	queryParams := []interface{}{userId}
	if q.Start != nil {
		queryStr += " AND end >= ?"
		queryParams = append(queryParams, q.Start.UTC())
	}
	if q.End != nil {
		queryStr += " AND start <= ?"
		queryParams = append(queryParams, q.End.UTC())
	}
	if q.ProjectId != nil {
		if q.IncludeSubProjects {
			queryStr += " AND (project_id = ? OR parent_project_id = ?)"
			queryParams = append(queryParams, q.ProjectId, q.ProjectId)
		} else {
			queryStr += " AND project_id = ?"
			queryParams = append(queryParams, q.ProjectId)
		}
	}
	if q.TodoId != nil {
		queryStr += " AND todo_id = ?"
		queryParams = append(queryParams, q.TodoId)
	}
	queryStr += " ORDER BY start, end IS NULL ASC"

	db, err := mysql.Open()
	if err != nil {
		return
	}
	defer db.Close()

	stmtOut, err := db.Prepare(queryStr)
	if err != nil {
		return
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query(queryParams...)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		p := Pomodoro{}
		err = rows.Scan(&p.Id, &p.Start, &p.End, &p.RemainingTime, &p.TodoId, &p.ProjectId, &p.ParentProjectId)
		if err != nil {
			return
		}
		pomodoros = append(pomodoros, p)
	}
	return
}
