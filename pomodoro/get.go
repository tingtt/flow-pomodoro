package pomodoro

import (
	"database/sql"
	"flow-pomodoros/mysql"
	"time"
)

func Get(userId uint64, id uint64) (p Pomodoro, notFound bool, err error) {
	db, err := mysql.Open()
	if err != nil {
		return
	}
	defer db.Close()

	stmtOut, err := db.Prepare("SELECT start, end, todo_id, project_id, parent_project_id FROM pomodoros WHERE user_id = ? AND id = ?")
	if err != nil {
		return
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query(userId, id)
	if err != nil {
		return
	}

	// TODO: uint64に対応
	var (
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
	err = rows.Scan(&start, &end, &todoId, &projectId, &rootProjectId)
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
