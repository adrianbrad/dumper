package psql

import (
	"database/sql"
	"dumper"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"sync"
)

type PayloadService struct {
	db             *sql.DB
	connectionInfo string
	mutex          sync.Mutex
}

func New(host, port, user, pass, name string) *PayloadService {
	psqlInfo := fmt.Sprintf(`
		host=%s
		port=%s 
		user=%s
		password=%s
		dbname=%s
		sslmode=disable`,
		host, port, user, pass, name)

	return &PayloadService{
		connectionInfo: psqlInfo,
	}
}

func (s *PayloadService) Open() (err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.db, err = sql.Open("postgres", s.connectionInfo)
	if err != nil {
		log.Errorf("Error while connecting to db, err :%s", err.Error())
		return
	}

	if err = s.db.Ping(); err != nil {
		log.Errorf("Error while pinging db, err: %s", err.Error())
		return
	}
	log.Info("Successfully connected to db")
	return
}

func (s *PayloadService) Close() error {
	return s.db.Close()
}

func (s *PayloadService) Write(p []byte) (n int, err error) {
	var payload dumper.Payload
	err = json.Unmarshal(p, &payload)
	if err != nil {
		return
	}

	if _, err = s.db.Exec(`
	INSERT INTO payloads(driver_id, longitude, latitude) 
	VALUES($1, $2, $3)`,
		payload.DriverID, payload.Longitude, payload.Latitude); err != nil {
		log.Errorf("Error while persisting payload to db, err: %s", err.Error())
		return
	}

	n = len(p)
	log.Info("Successfully wrote to db")
	return
}

func (s *PayloadService) Read(p []byte) (n int, err error) {
	var payload dumper.Payload

	s.mutex.Lock()
	if err = s.db.QueryRow(`
	SELECT driver_id, longitude, latitude
	FROM payloads
	ORDER BY created_at DESC
	LIMIT 1
	`).Scan(
		&payload.DriverID,
		&payload.Longitude,
		&payload.Latitude,
	); err != nil {
		log.Errorf("Error while querying payload from db, err: %s", err.Error())
		return
	}
	if _, err = s.db.Exec(`
	DELETE FROM payloads
	WHERE CTID = (
	    SELECT CTID
	    FROM payloads
	    ORDER BY created_at DESC
		LIMIT 1
	)
	`); err != nil {
		log.Errorf("Error while deleting payload from db, err: %s", err.Error())
		return
	}
	s.mutex.Unlock()

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return
	}
	n = copy(p, payloadBytes)
	return
}
