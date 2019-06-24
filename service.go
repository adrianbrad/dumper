package dumper

import "io"

type Opener interface {
	Open() (err error)
}

type Service interface {
	Opener
	io.ReadWriteCloser
}
