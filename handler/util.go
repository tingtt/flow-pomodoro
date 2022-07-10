package handler

import (
	"fmt"
	"strconv"
	"time"
)

func datetimeStrConv(str string) (t time.Time, err error) {
	// y-m-dTh:m:s or unix timestamp
	t, err1 := time.Parse("2006-1-2T15:4:5", str)
	if err1 == nil {
		return
	}
	t, err2 := time.Parse(time.RFC3339, str)
	if err2 == nil {
		return
	}
	u, err3 := strconv.ParseInt(str, 10, 64)
	if err3 == nil {
		t = time.Unix(u, 0)
		return
	}
	err = fmt.Errorf("\"%s\" is not a unix timestamp or string format \"2006-1-2T15:4:5\"", str)
	return
}
