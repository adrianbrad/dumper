package dumper

import "database/sql"

type Executor interface {
	Exec(query string, args ...interface{}) (result sql.Result, err error)
}
