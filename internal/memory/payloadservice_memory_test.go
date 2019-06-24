package memory

import (
	"dumper"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPayloadService_Write(t *testing.T) {
	s := New()
	_ = s.Open()
	p := dumper.Payload{
		DriverID:  1,
		Latitude:  2,
		Longitude: 3,
	}
	b, err := json.Marshal(p)
	assert.NoError(t, err)

	n, err := s.Write(b)
	assert.NoError(t, err)
	assert.Equal(t, len(b), n)

	assert.Equal(t, p, s.payloads[0])
}

func TestPayloadService_Read(t *testing.T) {
	s := New()
	_ = s.Open()
	p := dumper.Payload{
		DriverID:  1,
		Latitude:  2,
		Longitude: 3,
	}
	s.payloads = append(s.payloads, p)

	b := make([]byte, 100)
	n, err := s.Read(b)
	pBytes := b[:n]

	var receivedPayload dumper.Payload
	err = json.Unmarshal(pBytes, &receivedPayload)
	assert.NoError(t, err)

	assert.Equal(t, receivedPayload, p)
}

func TestPayloadService_Write_Race(t *testing.T) {
	s := New()
	_ = s.Open()
	p := dumper.Payload{
		DriverID:  1,
		Latitude:  2,
		Longitude: 3,
	}
	b, err := json.Marshal(p)
	assert.NoError(t, err)

	for i := 0; i <= 100; i++ {
		go func() {
			_, _ = s.Write(b)
		}()
	}
}

func TestPayloadService_Read_Race(t *testing.T) {
	s := New()
	_ = s.Open()
	s.payloads = append(s.payloads, dumper.Payload{})

	for i := 0; i <= 100; i++ {
		go func() {
			b := make([]byte, 100)
			_, _ = s.Read(b)
		}()
	}
}

func TestPayloadService_ReadWrite_Race(t *testing.T) {
	s := New()
	_ = s.Open()
	s.payloads = append(s.payloads, dumper.Payload{})

	p := dumper.Payload{
		DriverID:  1,
		Latitude:  2,
		Longitude: 3,
	}
	b, err := json.Marshal(p)
	assert.NoError(t, err)

	for i := 0; i <= 1000; i++ {
		go func() {
			b := make([]byte, 100)
			_, _ = s.Read(b)
		}()
		go func() {
			_, _ = s.Write(b)
		}()
	}
}
