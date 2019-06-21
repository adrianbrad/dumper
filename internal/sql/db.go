package sql

import "database/sql"

func A() {
	a, _ := sql.Open("", "")
	b, _ := a.Prepare("")
	_ = b.Close()
	_, _ = b.Exec("")
}
