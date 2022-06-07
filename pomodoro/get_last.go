package pomodoro

import (
	"flow-pomodoro/mysql"
)

func GetLast(userId uint64) (p Pomodoro, notFound bool, err error) {
	db, err := mysql.Open()
	if err != nil {
		return
	}
	defer db.Close()

	stmtOut, err := db.Prepare("SELECT id, start, end, remaining_time, todo_id, project_id, parent_project_id FROM logs WHERE user_id = ? ORDER BY end IS NULL DESC, start DESC, end DESC LIMIT 1")
	if err != nil {
		return
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query(userId)
	if err != nil {
		return
	}
	defer rows.Close()

	if !rows.Next() {
		// Not found
		notFound = true
		return
	}
	p = Pomodoro{}
	err = rows.Scan(&p.Id, &p.Start, &p.End, &p.RemainingTime, &p.TodoId, &p.ProjectId, &p.ParentProjectId)
	if err != nil {
		return
	}

	return
}
