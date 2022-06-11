package pomodoro

import (
	"flow-pomodoro/mysql"
	"time"
)

type PostStart struct {
	Start           time.Time `json:"start" validate:"required"`
	RemainingTime   *uint     `json:"remaining_time" validate:"omitempty"`
	TodoId          uint64    `query:"todo_id" json:"todo_id" validate:"required"`
	ProjectId       *uint64   `query:"project_id" json:"project_id" validate:"omitempty"`
	ParentProjectId *uint64   `query:"parent_project_id" json:"parent_project_id" validate:"omitempty"`
}

func Start(userId uint64, post PostStart, force bool) (p Pomodoro, notEnded bool, err error) {
	// Check ended
	old, notFound, err := GetLast(userId)
	if err != nil {
		return
	}
	if !notFound && old.End == nil {
		// Not ended
		notEnded = true
		if !force {
			return
		}
		// End last pomodoro
		var invalidTimeToEnd bool
		_, _, invalidTimeToEnd, err = End(userId, PostEnd{End: post.Start, TodoId: old.TodoId})
		if err != nil {
			return
		}
		if invalidTimeToEnd {
			_, err = Delete(userId, old.Id)
			if err != nil {
				return
			}
		}
	}

	// Insert
	db, err := mysql.Open()
	if err != nil {
		return
	}
	defer db.Close()
	stmtIns, err := db.Prepare("INSERT INTO logs (user_id, start, remaining_time, todo_id, project_id, parent_project_id) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return
	}
	defer stmtIns.Close()
	result, err := stmtIns.Exec(userId, post.Start.UTC(), post.RemainingTime, post.TodoId, post.ProjectId, post.ParentProjectId)
	if err != nil {
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		return
	}

	p.Id = uint64(id)
	p.Start = post.Start
	p.RemainingTime = post.RemainingTime
	p.ProjectId = post.ProjectId
	p.ParentProjectId = post.ParentProjectId

	return
}
