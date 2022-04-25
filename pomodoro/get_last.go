package pomodoro

import (
	"database/sql"
	"flow-pomodoro/mysql"
	"time"
)

func GetLast(userId uint64) (p Pomodoro, notFound bool, err error) {
	db, err := mysql.Open()
	if err != nil {
		return
	}
	defer db.Close()

	stmtOut, err := db.Prepare("SELECT id, start, end, todo_id, project_id, parent_project_id FROM logs WHERE user_id = ? ORDER BY end IS NULL DESC, start DESC, end DESC LIMIT 1")
	if err != nil {
		return
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query(userId)
	if err != nil {
		return
	}

	var (
		id              uint64
		start           time.Time
		end             sql.NullTime
		todoId          uint64
		projectId       sql.NullInt64
		parentProjectId sql.NullInt64
	)
	if !rows.Next() {
		// Not found
		notFound = true
		return
	}
	err = rows.Scan(&id, &start, &end, &todoId, &projectId, &parentProjectId)
	if err != nil {
		return Pomodoro{}, false, err
	}

	p.Id = id
	p.Start = start
	if end.Valid {
		p.End = &end.Time
	}
	p.TodoId = todoId
	if projectId.Valid {
		projectIdTmp := uint64(projectId.Int64)
		p.ProjectId = &projectIdTmp
	}
	if parentProjectId.Valid {
		parentProjectIdTmp := uint64(parentProjectId.Int64)
		p.ParentProjectId = &parentProjectIdTmp
	}

	return
}
