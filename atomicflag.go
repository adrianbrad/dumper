package dumper

import "sync/atomic"

type atomicFlag struct {
	flag int32
}

func (b *atomicFlag) Set(value bool) {
	if value {
		atomic.StoreInt32(&(b.flag), 1)
		return
	}
	atomic.StoreInt32(&(b.flag), 0)
}

func (b *atomicFlag) Get() bool {
	return atomic.LoadInt32(&(b.flag)) != 0
}
