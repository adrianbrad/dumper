package file

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"runtime"
	"strings"
	"sync"
	"testing"
)

var (
	initFileOnce     sync.Once
	initFileOnceFunc = func() {
		initFileOnce.Do(func() {
			_, err := os.Create("test.out")
			if err != nil {
				panic(err)
			}
			_, fileName, _, _ := runtime.Caller(0)
			dir := fileName[:strings.LastIndex(fileName, "/")]
			output = dir + "/test.out"
		})
	}
	output string
)

func TestPayloadService_Write(t *testing.T) {
	initFileOnceFunc()

	s := New(output)
	err := s.Open()
	assert.NoError(t, err)

	_, err = s.Write([]byte("test"))
	assert.NoError(t, err)
}

func TestPayloadService_Read(t *testing.T) {
	initFileOnceFunc()

	s := New(output)
	err := s.Open()
	assert.NoError(t, err)

	s.Write([]byte("hey"))

	p := make([]byte, 100)
	//n, err := s.Read(p)
	//assert.NoError(t, err)
	//assert.Equal(t, "hey", string(p[:n]))
	//
	//n, err = s.Read(p)
	//assert.Error(t, err)
	//assert.Equal(t, 0, n)

	for {
		_, err = s.Read(p)
		if err != nil {
			break
		}
		fmt.Print("wtf")
	}
}

func TestPayloadService_ReadWrite_Race(t *testing.T) {
	initFileOnceFunc()

	s := New(output)

	_ = s.Open()

	for i := 0; i <= 1000; i++ {
		go func() {
			b := make([]byte, 100)
			_, _ = s.Read(b)
		}()
		go func() {
			_, err := s.Write([]byte("hey"))
			assert.NoError(t, err)
		}()
	}
}
