package psql

import (
	"database/sql"
	"dumper"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

var (
	db            *sql.DB
	initDBOnce    sync.Once
	insertPayload = `
	INSERT INTO payloads(driver_id, latitude ,longitude)
	VALUES ($1, $2, $3);
	`
)

func openDB() {
	initDBOnce.Do(func() {
		psqlInfo := fmt.Sprintf(`
		host=%s
		port=%s 
		user=%s
		password=%s
		dbname=%s
		sslmode=disable`,
			"localhost", "5432", "admin", "admin", "dumper")
		d, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			panic(err)
		}
		db = d
	})
}

func truncate() {
	if _, err := db.Exec("TRUNCATE TABLE payloads"); err != nil {
		panic(err)
	}
}
func insertMock() {
	if _, err := db.Exec(insertPayload, 1, 100, 200); err != nil {
		panic(err)
	}
	if _, err := db.Exec(insertPayload, 2, 400, 123); err != nil {
		panic(err)
	}
	if _, err := db.Exec(insertPayload, 1, 5000, 1); err != nil {
		panic(err)
	}
}

func initService() *PayloadService {
	openDB()
	return New("localhost", "5432", "admin", "admin", "dumper")
}

func TestPayloadService_Open(t *testing.T) {
	s := initService()
	err := s.Open()
	assert.NoError(t, err)
}

func TestPayloadService_Read(t *testing.T) {
	s := initService()
	truncate()
	insertMock()
	_ = s.Open()

	b := make([]byte, 100)
	n, err := s.Read(b)
	assert.NoError(t, err)

	var p dumper.Payload
	err = json.Unmarshal(b[:n], &p)
	assert.NoError(t, err)

	expectedPayload := dumper.Payload{
		DriverID:  1,
		Latitude:  5000,
		Longitude: 1,
	}

	assert.Equal(t, expectedPayload, p)
}

func TestPayloadService_OpenCloseOpen(t *testing.T) {
	s := initService()
	err := s.Open()
	assert.NoError(t, err)
	err = s.Close()
	assert.NoError(t, err)
	err = s.Open()
	assert.NoError(t, err)
}

func TestPayloadService_Write(t *testing.T) {
	s := initService()
	truncate()
	insertMock()

	err := s.Open()
	assert.NoError(t, err)

	p := dumper.Payload{
		DriverID:  3,
		Latitude:  984,
		Longitude: 231,
	}

	b, err := json.Marshal(p)
	assert.NoError(t, err)

	n, err := s.Write(b)
	assert.Equal(t, len(b), n)
	assert.NoError(t, err)

	var insertedPayload dumper.Payload
	err = db.QueryRow(`
	SELECT driver_id, latitude, longitude
	FROM payloads
	ORDER BY created_at DESC
	LIMIT 1`).Scan(
		&insertedPayload.DriverID,
		&insertedPayload.Latitude,
		&insertedPayload.Longitude)
	assert.NoError(t, err)
	assert.Equal(t, insertedPayload, p)
}
