package pomodoro

import (
	"database/sql"
	"flow-pomodoro/mysql"
	"time"
)

func GetList(userId uint64, start time.Time, end time.Time) (pomodoros []Pomodoro, err error) {
	db, err := mysql.Open()
	if err != nil {
		return
	}
	defer db.Close()

	stmtOut, err := db.Prepare("SELECT id, start, end, todo_id, project_id, parent_project_id FROM logs WHERE user_id = ? AND (? BETWEEN start AND end OR ? BETWEEN start AND end OR start BETWEEN ? AND ?)")
	if err != nil {
		return
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query(userId, start, end, start, end)
	if err != nil {
		return
	}

	for rows.Next() {
		p := Pomodoro{}
		err = rows.Scan(&p.Id, &p.Start, &p.End, &p.TodoId, &p.ProjectId, &p.ParentProjectId)
		if err != nil {
			return
		}
		pomodoros = append(pomodoros, p)
	}

	return
}

func GetListTodo(userId uint64, start time.Time, end time.Time, todoId uint64) (pomodoros []Pomodoro, err error) {
	db, err := mysql.Open()
	if err != nil {
		return
	}
	defer db.Close()

	stmtOut, err := db.Prepare("SELECT id, start, end, project_id, parent_project_id FROM logs WHERE user_id = ? AND todo_id = ? AND (? BETWEEN start AND end OR ? BETWEEN start AND end OR start BETWEEN ? AND ?)")
	if err != nil {
		return
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query(userId, todoId, start, end, start, end)
	if err != nil {
		return
	}

	for rows.Next() {
		p := Pomodoro{TodoId: todoId}
		err = rows.Scan(&p.Id, &p.Start, &p.End, &p.ProjectId, &p.ParentProjectId)
		if err != nil {
			return
		}
		pomodoros = append(pomodoros, p)
	}

	return
}

func GetListProjectId(userId uint64, start time.Time, end time.Time, projectId uint64, includeSubProject bool) (pomodoros []Pomodoro, err error) {
	// Generate query
	queryStr := "SELECT id, start, end, todo_id, project_id, parent_project_id FROM logs WHERE user_id = ? AND (? BETWEEN start AND end OR ? BETWEEN start AND end OR start BETWEEN ? AND ?)"
	if includeSubProject {
		queryStr += " AND (project_id = ? OR parent_project_id = ?)"
	} else {
		queryStr += " AND project_id = ?"
	}

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

	var rows *sql.Rows
	if includeSubProject {
		rows, err = stmtOut.Query(userId, start, end, start, end, projectId, projectId)
	} else {
		rows, err = stmtOut.Query(userId, start, end, start, end, projectId)
	}
	if err != nil {
		return
	}

	for rows.Next() {
		p := Pomodoro{}
		err = rows.Scan(&p.Id, &p.Start, &p.End, &p.TodoId, &p.ProjectId, &p.ParentProjectId)
		if err != nil {
			return
		}
		pomodoros = append(pomodoros, p)
	}

	return
}
