package file

import (
	"fmt"
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
	return

}

func (s *PayloadService) Write(p []byte) (n int, err error) {
	s.mutex.Lock()
	n, err = fmt.Fprintln(s.f, string(p))
	s.lastLineLen = int64(n)
	s.mutex.Unlock()
	return
}

func (s *PayloadService) Read(p []byte) (n int, err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	input, err := ioutil.ReadFile(s.path)
	if err != nil {
		return
	}
	lines := strings.Split(string(input), "\n")
	linesLen := len(lines)
	if linesLen < 2 {
		err = fmt.Errorf("no lines to read")
		return
	}

	payload := lines[linesLen-2]
	output := strings.Join(lines[:linesLen-2], "\n")
	err = ioutil.WriteFile(s.path, []byte(output), 0644)
	if err != nil {
		return
	}

	n = copy(p, payload)
	return
}

func (s *PayloadService) Close() (err error) {
	return s.f.Close()
}
