package dbfilemem

import (
	"dumper/internal/file"
	"dumper/internal/http"
	"dumper/internal/memory"
	"dumper/internal/psql"
	"os"
)

func Run() {
	mem := memory.New()
	output := os.Getenv("OUT")
	f := file.New(output)
	host, port, user, pass, dbname := os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME")
	db := psql.New(host, port, user, pass, dbname)

	mem.Open()
	f.Open()
	db.Open()
	//
	//go func() {
	//	for {
	//		t := time.NewTicker(1 * time.Second)
	//		<-t.C
	//		db.Close()
	//	}
	//}()
	//fileMemDumper := dumper.New(2000, f, mem)
	//
	//dbFileDumper := dumper.New(2000, db, fileMemDumper)

	http.New(":8080", f)
}
