package memory

import (
	"dumper"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"sync"
)

type PayloadService struct {
	open     bool
	payloads []dumper.Payload
	mutex    sync.RWMutex
}

func New() *PayloadService {
	return &PayloadService{
		payloads: make([]dumper.Payload, 0),
	}
}

func (s *PayloadService) Open() (err error) {
	s.open = true
	log.Info("Successfully opened memory")
	return nil
}

func (s *PayloadService) Write(p []byte) (n int, err error) {
	if !s.open {
		err = fmt.Errorf("the payload service memory is closed")
		return
	}

	var payload dumper.Payload
	if err = json.Unmarshal(p, &payload); err != nil {
		return
	}

	s.mutex.Lock()
	s.payloads = append(s.payloads, payload)
	s.mutex.Unlock()

	log.Info("Successfully saved payload to file")
	n = len(p)
	return
}

func (s *PayloadService) Read(p []byte) (n int, err error) {
	if !s.open {
		err = fmt.Errorf("the payload service memory is closed")
		return
	}

	s.mutex.Lock()
	pLen := len(s.payloads)
	if pLen == 0 {
		err = fmt.Errorf("no elements in memory")
		return
	}
	var payload dumper.Payload

	// pop payload
	payload, s.payloads = s.payloads[pLen-1], s.payloads[:pLen-1]
	s.mutex.Unlock()

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return
	}

	n = copy(p, payloadBytes)
	log.Info("Successfully read payload from memory")
	return
}

func (s *PayloadService) Close() (err error) {
	s.open = false
	log.Info("Successfully closed memory")
	return
}
