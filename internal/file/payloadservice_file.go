package file

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

type PayloadService struct {
	f           *os.File
	path        string
	lastLineLen int64
	mutex       sync.Mutex
}

func New(path string) *PayloadService {
	return &PayloadService{
		path: path,
	}
}

func (s *PayloadService) Open() (err error) {
	s.f, err = os.OpenFile(s.path, os.O_APPEND|os.O_RDWR, os.ModeAppend)
	if err != nil {
		log.Errorf("Error while openening file, err: %s", err.Error())
		return
	}
	log.Info("Successfully opened file")
	return

}

func (s *PayloadService) Write(p []byte) (n int, err error) {
	s.mutex.Lock()
	n, err = fmt.Fprintln(s.f, string(p))
	if err != nil {
		log.Errorf("Error while writing to file, err: %s", err.Error())
		s.mutex.Unlock()
		return
	}

	s.lastLineLen = int64(n)
	s.mutex.Unlock()
	log.Info("Successfully saved to file")
	return
}

func (s *PayloadService) Read(p []byte) (n int, err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	input, err := ioutil.ReadFile(s.path)
	if err != nil {
		log.Errorf("Error while reading file, err: %s", err.Error())
		return
	}
	lines := strings.Split(string(input), "\n")
	linesLen := len(lines)
	if linesLen < 2 {
		err = fmt.Errorf("no lines to read")
		log.Error(err.Error())
		return
	}

	payload := lines[linesLen-2]
	output := strings.Join(lines[:linesLen-2], "\n")
	err = ioutil.WriteFile(s.path, []byte(output), 0644)
	if err != nil {
		log.Errorf("Error while writing to file, err: %s", err.Error())
		return
	}

	n = copy(p, payload)
	log.Info("Successfully read from file")
	return
}

func (s *PayloadService) Close() (err error) {
	log.Info("Closing file...")
	return s.f.Close()
}
