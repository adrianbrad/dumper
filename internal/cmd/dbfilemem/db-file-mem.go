package dbfilemem

import (
	"dumper"
	"dumper/internal/file"
	"dumper/internal/http"
	"dumper/internal/memory"
	"dumper/internal/psql"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

func Run() {
	output := os.Getenv("OUT")
	host, port, user, pass, dbname := os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME")
	srvport := os.Getenv("PORT")
	tic := os.Getenv("TIC")
	i, _ := strconv.Atoi(tic)

	os.Create(output)
	f := file.New(output)
	mem := memory.New()
	dbdump := psql.New(host, port, user, pass, dbname)

	fileMemDumper := dumper.New(i, f, mem)
	dbFileDumper := dumper.New(i, dbdump, fileMemDumper)

	f.Open()
	dbdump.Open()

	go func() {
		stopDB := time.NewTicker(1 * time.Second)

		select {
		case <-stopDB.C:
			log.Info("stopping db")
			dbdump.Close()

		}
	}()

	log.Error(http.New(srvport, dbFileDumper))
}
