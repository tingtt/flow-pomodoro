package pomodoro

import (
	"database/sql"
	"flow-pomodoros/mysql"
	"time"
)

func GetList(userId uint64, start time.Time, end time.Time) (pomodoros []Pomodoro, err error) {
	db, err := mysql.Open()
	if err != nil {
		return
	}
	defer db.Close()

	stmtOut, err := db.Prepare("SELECT id, start, end, todo_id, project_id, parent_project_id FROM pomodoros WHERE user_id = ? AND (? BETWEEN start AND end OR ? BETWEEN start AND end OR start BETWEEN ? AND ?)")
	if err != nil {
		return
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query(userId, start, end, start, end)
	if err != nil {
		return
	}

	for rows.Next() {
		var (
			id            uint64
			start         time.Time
			end           sql.NullTime
			todoId        uint64
			projectId     uint64
			rootProjectId uint64
		)
		err = rows.Scan(&id, &start, &end, &todoId, &projectId, &rootProjectId)
		if err != nil {
			return
		}

		p := Pomodoro{Id: id, Start: start, TodoId: todoId, ProjectId: projectId, ParentProjectId: rootProjectId}
		if end.Valid {
			p.End = &end.Time
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

	stmtOut, err := db.Prepare("SELECT id, start, end, project_id, parent_project_id FROM pomodoros WHERE user_id = ? AND todo_id = ? AND (? BETWEEN start AND end OR ? BETWEEN start AND end OR start BETWEEN ? AND ?)")
	if err != nil {
		return
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query(userId, todoId, start, end, start, end)
	if err != nil {
		return
	}

	for rows.Next() {
		var (
			id            uint64
			start         time.Time
			end           sql.NullTime
			projectId     uint64
			rootProjectId uint64
		)
		err = rows.Scan(&id, &start, &end, &projectId, &rootProjectId)
		if err != nil {
			return
		}

		p := Pomodoro{Id: id, Start: start, TodoId: todoId, ProjectId: projectId, ParentProjectId: rootProjectId}
		if end.Valid {
			p.End = &end.Time
		}

		pomodoros = append(pomodoros, p)
	}

	return
}

func GetListProjectId(userId uint64, start time.Time, end time.Time, projectId uint64, includeSubProject bool) (pomodoros []Pomodoro, err error) {
	// Generate query
	queryStr := "SELECT id, start, end, todo_id, project_id, parent_project_id FROM pomodoros WHERE user_id = ? AND (? BETWEEN start AND end OR ? BETWEEN start AND end OR start BETWEEN ? AND ?)"
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
		var (
			id            uint64
			start         time.Time
			end           sql.NullTime
			todoId        uint64
			projectId1    uint64
			rootProjectId uint64
		)
		err = rows.Scan(&id, &start, &end, &todoId, &projectId1, &rootProjectId)
		if err != nil {
			return
		}

		p := Pomodoro{Id: id, Start: start, TodoId: todoId, ProjectId: projectId1, ParentProjectId: rootProjectId}
		if end.Valid {
			p.End = &end.Time
		}

		pomodoros = append(pomodoros, p)
	}

	return
}
