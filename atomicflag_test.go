package dumper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAtomicFlag_SetTrue(t *testing.T) {
	var f atomicFlag
	f.Set(true)
	assert.True(t, f.Get())
}

func TestAtomicFlag_SetFalse(t *testing.T) {
	var f atomicFlag

	f.Set(false)
	assert.False(t, f.Get())
}

// run with -race flag
func TestAtomicFlag_Race(t *testing.T) {
	var f atomicFlag

	for i := 0; i <= 100; i++ {
		go func() {
			f.Set(!f.Get())
		}()
	}
}
