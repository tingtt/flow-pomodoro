package pomodoro

import (
	"flow-pomodoro/mysql"
)

func Get(userId uint64, id uint64) (p Pomodoro, notFound bool, err error) {
	db, err := mysql.Open()
	if err != nil {
		return
	}
	defer db.Close()

	stmtOut, err := db.Prepare("SELECT start, end, todo_id, project_id, parent_project_id FROM logs WHERE user_id = ? AND id = ?")
	if err != nil {
		return
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query(userId, id)
	if err != nil {
		return
	}
	defer rows.Close()

	if !rows.Next() {
		// Not found
		notFound = true
		return
	}
	err = rows.Scan(&p.Start, &p.End, &p.TodoId, &p.ProjectId, &p.ParentProjectId)
	if err != nil {
		return Pomodoro{}, false, err
	}

	p.Id = id
	return
}
