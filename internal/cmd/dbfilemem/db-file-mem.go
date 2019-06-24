package dbfilemem

import (
	"dumper"
	"dumper/internal/file"
	"dumper/internal/http"
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
	f := file.New(output)
	db := psql.New(host, port, user, pass, dbname)

	f.Open()
	db.Open()

	go func() {
		stopDB := time.NewTicker(1 * time.Second)

		//for {
		select {
		case <-stopDB.C:
			log.Info("stopping db")
			db.Close()

		}
		//}
	}()
	dbFileDumper := dumper.New(i, db, f)

	log.Error(http.New(srvport, dbFileDumper))
}
