package pomodoro

import (
	"flow-pomodoro/mysql"
	"time"
)

type PostEnd struct {
	End           time.Time `json:"end" validate:"required"`
	RemainingTime *uint     `json:"remaining_time" validate:"omitempty"`
	TodoId        uint64    `query:"todo_id" json:"todo_id" validate:"required"`
}

func End(userId uint64, post PostEnd) (p Pomodoro, notStarted bool, invalidTime bool, err error) {
	// Check started
	p, notFound, err := GetLast(userId)
	if err != nil {
		return
	}
	if notFound || p.End != nil || p.TodoId != post.TodoId {
		// Not started
		notStarted = true
		return
	}
	if post.End.Before(p.Start) {
		invalidTime = true
		return
	}

	// Update
	db, err := mysql.Open()
	if err != nil {
		return
	}
	defer db.Close()
	stmtIns, err := db.Prepare("UPDATE logs SET end  = ?, remaining_time = ? WHERE user_id = ? AND id = ?")
	if err != nil {
		return
	}
	defer stmtIns.Close()
	_, err = stmtIns.Exec(post.End.UTC(), post.RemainingTime, userId, p.Id)
	if err != nil {
		return
	}

	p.End = &post.End
	return
}
