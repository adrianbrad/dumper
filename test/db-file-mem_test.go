package test

import (
	"database/sql"
	"dumper"
	"dumper/internal/file"
	"dumper/internal/memory"
	"dumper/internal/psql"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"
)

func TestFileFallback(t *testing.T) {
	output := os.Getenv("OUT")
	host, port, user, pass, dbname := os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME")

	f := file.New(output)
	db := psql.New(host, port, user, pass, dbname)

	os.Create("test.out")

	dbFileDumper := dumper.New(1000, db, f)

	db.Open()
	db.Close()

	f.Open()

	p := []byte(`{"driver_id":1, "latitude":1, "longitude":2`)
	dbFileDumper.Write(p)

	f.Close()

	b, err := ioutil.ReadFile("./test.out")
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(b), "\n")

	assert.Equal(t, string(p), lines[0])
}

func TestFileFallbackThenWriteToDB(t *testing.T) {
	output := os.Getenv("OUT")
	host, port, user, pass, dbname := os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME")
	db, err := sql.Open("postgres", fmt.Sprintf(`
		host=%s
		port=%s 
		user=%s
		password=%s
		dbname=%s
		sslmode=disable`,
		host, port, user, pass, dbname))
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	db.Exec("truncate table payloads")

	f := file.New(output)
	dbdump := psql.New(host, port, user, pass, dbname)

	os.Create("test.out")

	dbFileDumper := dumper.New(100, dbdump, f)

	dbdump.Open()
	dbdump.Close()

	f.Open()

	p := []byte(`{"driver_id":1, "latitude":1, "longitude":1}`)
	p2 := []byte(`{"driver_id":2, "latitude":2, "longitude":2}`)
	dbFileDumper.Write(p)
	dbFileDumper.Write(p2)

	time.Sleep(500 * time.Millisecond)
	r, err := db.Query("select driver_id, latitude, longitude from payloads order by created_at desc")
	if err != nil {
		panic(err)
	}

	payload := dumper.Payload{}
	var i int64 = 0

	for r.Next() {
		i++
		err = r.Scan(&payload.DriverID, &payload.Latitude, &payload.Longitude)
		if err != nil {
			panic(err)
		}
		assert.Equal(t, i, payload.DriverID)
	}

	assert.Equal(t, int64(2), i)
}

func TestMemoryFallback(t *testing.T) {
	output := os.Getenv("OUT")
	host, port, user, pass, dbname := os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME")

	os.Create("test.out")
	f := file.New(output)
	dbdump := psql.New(host, port, user, pass, dbname)
	mem := memory.New()

	fileMemDumper := dumper.New(500, f, mem)
	dbFileDumper := dumper.New(500, dbdump, fileMemDumper)

	dbdump.Open()
	dbdump.Close()

	f.Open()
	f.Close()

	mem.Open()

	p := []byte(`{"driver_id":1,"latitude":1,"longitude":1}`)
	p2 := []byte(`{"driver_id":2,"latitude":2,"longitude":2}`)
	dbFileDumper.Write(p)
	dbFileDumper.Write(p2)

	pMem := make([]byte, 100)
	n, err := mem.Read(pMem)
	if err != nil {
		panic(err)
	}
	var payload dumper.Payload
	err = json.Unmarshal(pMem[:n], &payload)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, payload.DriverID, int64(2))
	assert.Equal(t, payload.Longitude, float64(2))
	assert.Equal(t, payload.Latitude, float64(2))

	n, err = mem.Read(pMem)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(pMem[:n], &payload)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, payload.DriverID, int64(1))
	assert.Equal(t, payload.Longitude, float64(1))
	assert.Equal(t, payload.Latitude, float64(1))
}

func TestMemoryThenWriteToFileThenWriteToDB(t *testing.T) {
	output := os.Getenv("OUT")
	host, port, user, pass, dbname := os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME")
	db, err := sql.Open("postgres", fmt.Sprintf(`
		host=%s
		port=%s 
		user=%s
		password=%s
		dbname=%s
		sslmode=disable`,
		host, port, user, pass, dbname))
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	db.Exec("truncate table payloads")

	os.Create(output)
	f := file.New(output)
	dbdump := psql.New(host, port, user, pass, dbname)
	mem := memory.New()

	fileMemDumper := dumper.New(300, f, mem)
	dbFileDumper := dumper.New(500, dbdump, fileMemDumper)

	dbdump.Open()
	dbdump.Close()

	f.Open()
	f.Close()

	mem.Open()

	p := []byte(`{"driver_id":1,"latitude":1,"longitude":1}`)
	p2 := []byte(`{"driver_id":2,"latitude":2,"longitude":2}`)
	dbFileDumper.Write(p)
	dbFileDumper.Write(p2)

	time.Sleep(5000 * time.Millisecond)
	r, err := db.Query("select driver_id, latitude, longitude from payloads order by created_at asc")
	if err != nil {
		panic(err)
	}

	payload := dumper.Payload{}
	var i int64 = 0

	for r.Next() {
		i++
		err = r.Scan(&payload.DriverID, &payload.Latitude, &payload.Longitude)
		if err != nil {
			panic(err)
		}
		assert.Equal(t, i, payload.DriverID)
	}

	assert.Equal(t, int64(2), i)
}
