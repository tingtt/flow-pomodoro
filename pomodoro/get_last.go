package pomodoro

import (
	"database/sql"
	"flow-pomodoros/mysql"
	"time"
)

func GetLast(userId uint64) (p Pomodoro, notFound bool, err error) {
	db, err := mysql.Open()
	if err != nil {
		return
	}
	defer db.Close()

	stmtOut, err := db.Prepare("SELECT id, start, end, todo_id, project_id, parent_project_id FROM pomodoros WHERE user_id = ? ORDER BY start DESC, end IS NULL DESC, end DESC LIMIT 1")
	if err != nil {
		return
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query(userId)
	if err != nil {
		return
	}

	var (
		id            uint64
		start         time.Time
		end           sql.NullTime
		todoId        uint64
		projectId     uint64
		rootProjectId uint64
	)
	if !rows.Next() {
		// Not found
		notFound = true
		return
	}
	err = rows.Scan(&id, &start, &end, &todoId, &projectId, &rootProjectId)
	if err != nil {
		return Pomodoro{}, false, err
	}

	p.Id = id
	p.Start = start
	if end.Valid {
		p.End = &end.Time
	}
	p.ProjectId = projectId
	p.ParentProjectId = rootProjectId

	return
}
