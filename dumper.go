package dumper

type Dumper interface {
	Dump(p Payload) (err error)
}
